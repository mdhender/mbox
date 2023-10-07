package newsgroup

import (
	"bytes"
	"github.com/mdhender/mbox/internal/chunk"
	"strings"
	"unicode"
)

// AddToCorpus add the words from a Chunk into the Corpus.
func (ng *NewsGroup) AddToCorpus(ch *chunk.Chunk) {
	for _, line := range ch.Body {
		// split on any non-letter/non-number rune.
		for _, word := range bytes.FieldsFunc(line, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		}) {
			word := strings.ToLower(string(word))
			if _, ok := ng.Corpus.StopWords[word]; ok {
				// filter out the stop word
				continue
			}
			ng.Corpus.Words[word] = ng.Corpus.Words[word] + 1
		}
	}
}
