package chunk

import (
	"bytes"
	stemmer "github.com/agonopol/go-stem"
	"log"
	"os"
	"regexp"
	"time"
	"unicode"
)

type Chunk struct {
	Line   int
	From   []byte
	Header [][]byte
	Body   [][]byte
}

// Chunks is evil. It splits the input into chunks and pre-processes the input, too.
func Chunks(path string) ([]*Chunk, error) {
	started := time.Now()

	input, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	log.Printf("[chunk] read %d bytes in %v\n", len(input), time.Now().Sub(started))

	var chunks []*Chunk

	// every message starts with a blank line followed by a From with ID.
	som := regexp.MustCompile("^From -?[0-9]+$")
	log.Printf("[chunk] completed compile in %v\n", time.Now().Sub(started))

	// split into lines and trim any trailing spaces
	lines := bytes.Split(input, []byte{'\n'})
	log.Printf("[chunk] completed split   in %v\n", time.Now().Sub(started))
	for i := 0; i < len(lines); i++ {
		lines[i] = bytes.TrimRight(lines[i], " \r\t")
	}
	log.Printf("[chunk] completed trim    in %v\n", time.Now().Sub(started))

	for n := 0; n < len(lines); n++ {
		line := lines[n]
		if som.Find(line) != nil {
			ch := &Chunk{
				Line: n + 1,
				From: line,
			}
			for n = n + 1; n < len(lines); n++ {
				line = lines[n]
				if len(line) == 0 {
					break
				} else if len(line) == 1 && (line[0] == ' ' || line[0] == '\t') {
					break
				}

				// header should never have tabs in it
				for i, ch := range line {
					if ch == '\t' {
						line[i] = ' '
					}
				}

				// pre-processing hacks
				if bytes.HasPrefix(line, []byte("Date: ")) {
					if bytes.Equal(line, []byte("Date: 11 Sep 93 12:58:28 -500")) {
						line = []byte("Date: 11 Sep 93 12:58:28 -0500")
					} else if bytes.Equal(line, []byte("Date: 11 Sep 93 23:10:45 -500")) {
						line = []byte("Date: 11 Sep 93 23:10:45 -0500")
					} else if bytes.Equal(line, []byte("Date: Wed, 12 Oct 1994 09:35:51 Central")) {
						line = []byte("Date: Wed, 12 Oct 1994 09:35:51 CST")
					} else if bytes.Equal(line, []byte("Date: Thu, 02 Dec 93 19:50:54 est")) {
						line = []byte("Date: Thu, 02 Dec 93 19:50:54 EST")
					} else if bytes.Equal(line, []byte("Date: Tue, 15 Jun 93 15:10:37 T-1")) {
						line = []byte("Date: Tue, 15 Jun 93 15:10:37 -0100")
					}
				} else if bytes.HasPrefix(line, []byte("References: ")) {
					if bytes.Equal(line, []byte("References: <")) {
						line = []byte("References: <missing-reference-id>")
					} else if bytes.Equal(line, []byte("References: C0GzED.A2u@news.cso.uiuc.edu> <1829@idacrd.UUCP> <1ii5rfINNc2q@darkstar.UCSC.EDU")) {
						line = []byte("References: <C0GzED.A2u@news.cso.uiuc.edu> <1829@idacrd.UUCP> <1ii5rfINNc2q@darkstar.UCSC.EDU>")
					} else if bytes.Equal(line, []byte("References: RSI Customer Service")) {
						line = []byte("References: <RSI-Customer-Service>")
					} else if bytes.Equal(line, []byte("References: <1991Apr13.030312.7999@vax1.tcd.ie}")) {
						line = []byte("References: <1991Apr13.030312.7999@vax1.tcd.ie>")
					} else if bytes.Equal(line, []byte("References: <1991Nov12.183857.24316@newcastle.ac.uk> <1991Nov18.011915.40")) {
						line = []byte("References: <1991Nov12.183857.24316@newcastle.ac.uk> <1991Nov18.011915.408@bradley.bradley.edu>")
					} else if bytes.Equal(line, []byte("References: <1992Mar21.004047.17322@erg.sri.com>> <18182@ector.cs.purdue.edu> <1992Mar21.213430.8671@daimi.aau.dk")) {
						line = []byte("References: <1992Mar21.004047.17322@erg.sri.com> <18182@ector.cs.purdue.edu> <1992Mar21.213430.8671@daimi.aau.dk>")
					} else if bytes.Equal(line, []byte("References: <1993Feb1.162305.16901@magnus.acs.ohio-state.edu> <1kjon1INN81d@bre")) {
						line = []byte("References: <1993Feb1.162305.16901@magnus.acs.ohio-state.edu> <1kjon1INN81d@bredbeddle.cs.purdue.edu>")
					} else if bytes.Equal(line, []byte("References: <8fJ=SMe00WBLE7En4P@andrew.cmu.edu> <21390@ucdavis.ucdavis.edu> <8f")) {
						line = []byte("References: <8fJ=SMe00WBLE7En4P@andrew.cmu.edu> <21390@ucdavis.ucdavis.edu> <invalid-reference-id>")
					} else if bytes.Equal(line, []byte("References: <C1tyDE.EI9@inews.Intel.COM> <16B69C2D4.X049RH@tamvm1.tamu.edu> <19")) {
						line = []byte("References: <C1tyDE.EI9@inews.Intel.COM> <16B69C2D4.X049RH@tamvm1.tamu.edu> <1993Feb4.044100.17009@midway.uchicago.edu>")
					}
				}

				if len(ch.Header) != 0 && line[0] == ' ' {
					ch.Header[len(ch.Header)-1] = append(ch.Header[len(ch.Header)-1], line...)
				} else {
					ch.Header = append(ch.Header, line)
				}
			}
			for n = n + 1; n < len(lines) && !endOfMessage(lines[n:]); n++ {
				line = lines[n]

				ch.Body = append(ch.Body, line)
			}
			chunks = append(chunks, ch)
		}
	}

	log.Printf("[chunk] completed chunks  in %v\n", time.Now().Sub(started))
	return chunks, nil
}

func (ch *Chunk) Words(stopWords map[string]bool) map[string]int {
	words := make(map[string]int)

	// split on any non-letter/non-number rune.
	for _, line := range ch.Body {
		for _, token := range Tokenize(line, stopWords) {
			word := string(token)
			words[word] = words[word] + 1
		}
		//for _, word := range bytes.FieldsFunc(line, func(r rune) bool {
		//	return !unicode.IsGraphic(r) || unicode.IsSpace(r) || r == '"' || r == '(' || r == ')' || r == ',' || r == '.'
		//}) {
		//	word := bytes.ToLower(word)
		//	if isOnlyLetters(word) { // filter out words that contain non-letters
		//		if !stopWords[string(word)] { // filter out stop-words
		//			// convert to stem word
		//			fjord := string(stemmer.Stem(word))
		//			words[fjord] = words[fjord] + 1
		//		}
		//	}
		//}
	}

	return words
}

func endOfMessage(lines [][]byte) bool {
	if len(lines) == 0 {
		return true
	} else if len(lines) == 2 && len(lines[0]) == 0 && len(lines[1]) == 0 {
		return true
	} else if len(lines) > 2 && len(lines[0]) == 0 && len(lines[1]) == 0 && bytes.HasPrefix(lines[2], []byte("From ")) {
		return true
	}
	return false
}

func isOnlyLetters(s []byte) bool {
	for _, ch := range s {
		if !('a' <= ch && ch <= 'z') {
			return false
		}
	}
	return true
}

func Tokenize(line []byte, stopWords map[string]bool) [][]byte {
	var tokens [][]byte

	for _, word := range bytes.FieldsFunc(line, func(r rune) bool {
		if unicode.IsSpace(r) {
			return true
		} else if unicode.IsLetter(r) {
			return false
		} else if unicode.IsDigit(r) {
			return false
		} else if r == '\t' || r == '"' || r == '(' || r == ')' || r == ',' || r == '.' {
			return true
		}
		return true
	}) {
		word := bytes.ToLower(word)
		if isOnlyLetters(word) { // filter out words that contain non-letters
			if len(word) > 3 { // avoid short words
				if !stopWords[string(word)] { // filter out stop-words
					// convert to stem word
					tokens = append(tokens, stemmer.Stem(word))
				}
			}
		}
	}

	return tokens
}
