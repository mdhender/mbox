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

	defer func(started time.Time) {
		log.Printf("[mbox] completed in %v\n", time.Now().Sub(started))
	}(time.Now())

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

	linesChecked, linesMatched := 0, 0
	for _, msg := range msgs {
		if err := msg.Parse(); err != nil {
			log.Printf("[mbox] message %d: %v\n", msg.Start, err)
			continue
		}
		// sanity check
		if msg.Header.Lines != 0 {
			linesChecked++
			if msg.Header.Lines != len(msg.Body.Text) {
				// log.Printf("[mbox] message %8d: body: lines: want %6d: got %6d\n", msg.Start, msg.Header.Lines, len(msg.Body.Text))
			} else {
				linesMatched++
			}
		}
	}
	log.Printf("[mbox] checked %8d messages; lines matched on %8d, missed on %8d\n", linesChecked, linesMatched, linesChecked-linesMatched)
	if len(msgs) != len(mbox) {
		log.Printf("[mbox] expected %8d messages: got %8d\n", len(msgs), len(mbox))
	}
	// link messages
	for _, msg := range msgs {
		if err := msg.Header.LinkReferences(mbox); err != nil {
			log.Printf("[mbox] message %d: %v\n", msg.Start, err)
			continue
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

	log.Fatalln(http.ListenAndServe(":8080", a.router))
}
