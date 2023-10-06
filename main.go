package main

import (
	"flag"
	"github.com/matryer/way"
	"log"
	"net/http"
	"sort"
	"time"
)

func main() {
	showHeaders := false
	flag.BoolVar(&showHeaders, "show-headers", showHeaders, "show headers")
	flag.Parse()

	started := time.Now()
	defer func(started time.Time) {
		log.Printf("[mbox] completed in %v\n", time.Now().Sub(started))
	}(started)

	lines, err := read("rec.games.pbm.mbox")
	if err != nil {
		log.Fatalf("[mbox] read mbox: %v\n", err)
	}
	log.Printf("[mbox] read %d lines\n", len(lines))

	lines = preprocess(lines)
	log.Printf("[mbox] pre-processed lines\n")

	msgs, err := split(lines)
	if err != nil {
		for _, msg := range msgs {
			if msg.Error != nil {
				log.Printf("[mbox] msg %d: %v\n", msg.Start, msg.Error)
			}
		}
		log.Fatalf("[mbox] split: %v\n", err)
	}
	log.Printf("[mbox] read %d messages\n", len(msgs))

	duplicateIds := 0
	for _, msg := range msgs {
		if err := msg.Parse(); err != nil {
			log.Printf("[mbox] message %d: %v\n", msg.Start, err)
			continue
		}
		// sanity check
		if mbox[msg.Header.Id] != nil {
			duplicateIds++
			log.Printf("[mbox] duplicate message %q\n", msg.Header.Id)
		}

		// set the message id to the header id
		msg.Id = msg.Header.Id
		mbox[msg.Id] = msg
		msg.Spam = spam[msg.Id]
		msg.Struck = struck[msg.Id]
	}
	if len(msgs) != len(mbox) {
		log.Printf("[mbox] expected %8d messages: got %8d\n", len(msgs), len(mbox))
	}
	if duplicateIds != 0 {
		log.Fatalf("[mbox] found %d duplicate message ids\n", duplicateIds)
	}

	// link messages forward and backwards
	for _, msg := range msgs {
		if err := msg.Header.LinkReferences(mbox); err != nil {
			log.Printf("[mbox] message %d: %v\n", msg.Start, err)
		}
		for _, ref := range msg.Header.References.Messages {
			msg.References = append(msg.References, ref)
			ref.ReferencedBy = append(ref.ReferencedBy, msg)
		}
	}

	headers, err := CollectHeaderKeys(msgs)
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

	a := &App{
		router:    way.NewRouter(),
		templates: "../templates",
	}
	a.router.HandleFunc("GET", "/messages/:id", a.handleMessage)

	a.messages.byId = mbox
	a.messages.byLine = make(map[int]*Message)
	for _, msg := range msgs {
		a.messages.byLine[msg.Start] = msg
	}

	for _, msg := range msgs {
		if spam[msg.Header.Id] || struck[msg.Header.Id] {
			continue
		}
		if msg.Header.From == "HGHFGDS <fhfgfgg@gmail.com>" {
			log.Printf("[spam] %q\n", msg.Header.Id)
		} else if msg.Header.From == "iwcwatches5@gmail.com" {
			log.Printf("[spam] %q\n", msg.Header.Id)
		}
	}

	log.Printf("[mbox] completed prep in %v\n", time.Now().Sub(started))

	log.Fatalln(http.ListenAndServe(":8080", a.router))
}
