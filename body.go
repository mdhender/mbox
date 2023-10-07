package main

import (
	"strings"
	"unicode"
)

type Body struct {
	Text  [][]byte
	Value string
}

func (b *Body) Parse(spam, hide bool) error {
	if spam || hide {
		b.Text = [][]byte{[]byte("Message content has been removed.")}
	}

	sb := strings.Builder{}
	for _, line := range b.Text {
		sb.Write(line)
		sb.WriteByte('\n')
	}
	b.Value = sb.String()

	return nil
}

// Words returns all the words in the body as a slice of strings
func (b *Body) Words() map[string]int {
	words := make(map[string]int)
	// split on any non-letter/non-number rune.
	for _, word := range strings.FieldsFunc(b.Value, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	}) {
		word := strings.ToLower(word)
		if _, ok := stopwords[word]; ok {
			// filter out the stop word
			continue
		}
		words[word] = words[word] + 1
	}
	return words
}
