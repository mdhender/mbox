// Package mbox implements a mailbox store.
package mbox

import (
	"fmt"
	"log"
	"sort"
)

type MailBox struct {
	Messages []*Message
	ById     map[string]*Message
	ByLine   map[int]*Message
	Corpus   map[string]int
}

func New(input [][]byte, showHeaders bool) (*MailBox, error) {
	mb := &MailBox{
		ById:   make(map[string]*Message),
		ByLine: make(map[int]*Message),
		Corpus: make(map[string]int),
	}

	// parse the raw lines into separate messages
	mb.Parse(input)
	log.Printf("[mbox] read %d messages\n", len(mb.Messages))

	// check for any errors from that parsing
	var err error
	for _, msg := range mb.Messages {
		if msg.Error != nil {
			err = fmt.Errorf("mbox: parse")
			log.Printf("[mbox] msg %d: %v\n", msg.Start, msg.Error)
		}
	}
	if err != nil {
		return nil, err
	}

	// parse the header
	mb.ParseHeaders()

	// parse more of the message
	duplicateIds := 0
	for _, msg := range mb.Messages {
		if err := msg.Parse(); err != nil {
			log.Printf("[mbox] message %d: %v\n", msg.Start, err)
			continue
		}
		// sanity check
		if mb.ById[msg.Header.Id] != nil {
			duplicateIds++
			log.Printf("[mbox] duplicate message %q\n", msg.Header.Id)
		}

		// set the message id to the header id
		msg.Id = msg.Header.Id
		mb.ById[msg.Id] = msg
	}
	if len(mb.ById) != len(mb.Messages) {
		log.Printf("[mbox] expected %8d messages: got %8d\n", len(mb.Messages), len(mb.ById))
	}
	if duplicateIds != 0 {
		log.Fatalf("[mbox] found %d duplicate message ids\n", duplicateIds)
	}

	// index by lines
	mb.ByLine = make(map[int]*Message)
	for _, msg := range mb.Messages {
		mb.ByLine[msg.Start] = msg
	}

	headers, err := collectHeaderKeys(mb.Messages)
	if err != nil {
		log.Fatalf("[mbox] collectHeaders: %v\n", err)
	} else if showHeaders {
		var text []string
		for k := range headers {
			text = append(text, k)
		}
		sort.Strings(text)
		for _, t := range text {
			log.Printf("[mbox] header %-35s == %8d\n", t, headers[t])
		}
	}
	log.Printf("[mbox] read %8d header types\n", len(headers))

	return mb, nil
}

func (mb *MailBox) LinkMessages() {
	// link messages forward and backwards
	for _, msg := range mb.Messages {
		if err := msg.Header.LinkReferences(mb.ById); err != nil {
			log.Printf("[mbox] message %d: %v\n", msg.Start, err)
		}
		for _, ref := range msg.Header.References.Messages {
			msg.References = append(msg.References, ref)
			ref.ReferencedBy = append(ref.ReferencedBy, msg)
		}
	}
}

func (mb *MailBox) ParseHeaders() {
	for _, msg := range mb.Messages {
		msg.Error = msg.Header.Parse()
		msg.Id = msg.Header.Id
		if mb.ById[msg.Id] != nil {
			log.Fatal("duplicate message id %q", msg.Id)
		}
		msg.Subject = msg.Header.Subject
		if msg.Spam = spam[msg.Header.Id]; msg.Spam {
			msg.Subject = "Message has been flagged as spam."
		}
		if msg.Struck = struck[msg.Header.Id]; msg.Struck {
			msg.Subject = "Message has been taken down due to request."
		}
	}

}
