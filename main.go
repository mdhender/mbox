package main

import (
	"bytes"
	"flag"
	"github.com/mdhender/mbox/internal/app"
	"github.com/mdhender/mbox/internal/stores/mbox"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"
)

func main() {
	doCorpus, showHeaders, flagSpam, flagStruck := false, false, false, false
	flag.BoolVar(&doCorpus, "corpus", doCorpus, "create corpus")
	flag.BoolVar(&flagSpam, "flag-spam", flagSpam, "show suspected spam headers")
	flag.BoolVar(&flagStruck, "flag-struck", flagStruck, "show suspected struct headers")
	flag.BoolVar(&showHeaders, "show-headers", showHeaders, "show headers")
	flag.Parse()

	started := time.Now()
	defer func(started time.Time) {
		log.Printf("[mbox] completed in %v\n", time.Now().Sub(started))
	}(started)

	input, err := os.ReadFile("rec.games.pbm.mbox")
	if err != nil {
		log.Fatalf("[mbox] read mbox: %v\n", err)
	}
	log.Printf("[mbox] read %d bytes in %v\n", len(input), time.Now().Sub(started))
	// split into lines and trim any carriage-returns
	lines := bytes.Split(input, []byte{'\n'})
	for i := 0; i < len(lines); i++ {
		lines[i] = bytes.TrimRight(lines[i], "\r")
	}
	log.Printf("[mbox] completed split        in %v\n", time.Now().Sub(started))
	lines = preProcess(lines)
	log.Printf("[mbox] completed pre-process  in %v\n", time.Now().Sub(started))

	if doCorpus {
		corpus := mkCorpus(input)
		log.Printf("[mbox] created corpus  in %v (%d)\n", time.Now().Sub(started), len(corpus))
	}

	box, err := mbox.New(lines, showHeaders)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[mbox] completed load  in %v\n", time.Now().Sub(started))

	if flagSpam || flagStruck {
		if flagSpam {
			box.FlagSpam()
		}
		if flagStruck {
			box.FlagStruck()
		}
		log.Printf("[mbox] completed flags in %v\n", time.Now().Sub(started))
		os.Exit(2)
	}

	box.LinkMessages()
	log.Printf("[mbox] linked messages in %v\n", time.Now().Sub(started))

	//box.MakeCorpus()
	log.Printf("[mbox] created corpus  in %v\n", time.Now().Sub(started))

	a, err := app.New(box)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatalln(http.ListenAndServe(":8080", a.Router))
}

// mkCorpus returns all the words in the body as a slice of strings
func mkCorpus(input []byte) map[string]int {
	words := make(map[string]int)
	// split on any non-letter/non-number rune.
	for _, word := range bytes.FieldsFunc(input, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	}) {
		word := strings.ToLower(string(word))
		if _, ok := stopwords[word]; ok {
			// filter out the stop word
			continue
		}
		words[word] = words[word] + 1
	}
	return words
}
