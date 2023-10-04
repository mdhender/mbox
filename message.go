package main

import (
	"bytes"
)

type Message struct {
	Start  int      // line number in mbox file
	Raw    [][]byte // the line starting the message
	RawId  string   // ID from the "From " block header
	Header *Header
	Body   *Body
	Spam   bool
	Struck bool  // struck for copyright or ownership
	Error  error // any error parsing the message
}

func (m *Message) Parse() error {
	// populate header and body with raw lines
	for _, line := range m.Raw {
		if m.Body != nil {
			m.Body.Text = append(m.Body.Text, line)
		} else if m.Header == nil && bytes.HasPrefix(line, []byte("From ")) {
			m.RawId = string(line[5:])
			m.Struck = struck[m.RawId]
			m.Header = &Header{}
		} else if len(line) == 0 {
			m.Body = &Body{}
		} else if len(m.Header.Text) != 0 && (line[0] == ' ' || line[0] == '\t') {
			m.Header.Text[len(m.Header.Text)-1] = append(m.Header.Text[len(m.Header.Text)-1], '.')
			m.Header.Text[len(m.Header.Text)-1] = append(m.Header.Text[len(m.Header.Text)-1], line[1:]...)
		} else {
			m.Header.Text = append(m.Header.Text, line)
		}
	}

	// buggy hack - always slice off the final two lines of the body
	if m.Body != nil {
		// remove the last line if it isn't empty
		for i := 0; i < 2; i++ {
			if len(m.Body.Text) != 0 && len(m.Body.Text[len(m.Body.Text)-1]) == 0 {
				m.Body.Text = m.Body.Text[:len(m.Body.Text)-1]
			}
		}
	}

	// parse header
	if err := m.Header.Parse(spam, struck); err != nil {
		return err
	}
	mbox[m.Header.Id] = m
	m.Spam = spam[m.Header.Id]
	m.Struck = struck[m.Header.Id]

	// parse body
	if err := m.Body.Parse(m.Spam, m.Struck); err != nil {
		return err
	}

	return nil
}
