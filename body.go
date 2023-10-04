package main

import "strings"

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
