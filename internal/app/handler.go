package app

import (
	"github.com/matryer/way"
	"github.com/mdhender/mbox/internal/stores/mbox"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (a *App) handleCorpus(w http.ResponseWriter, r *http.Request) {
	log.Printf("[corpus] dump %d\n", len(a.MailBox.Corpus))
	payload := struct {
		Corpus map[string]int
	}{
		Corpus: a.MailBox.Corpus,
	}
	a.render(w, r, payload, "layout", "corpus")
}

func (a *App) handleIndex() http.HandlerFunc {
	type Bucket struct {
		Period   string
		Messages []*mbox.Message
	}
	var payload struct {
		ArticleCount int
		From         string
		Through      string
		ByYear       map[string]*Bucket
	}
	payload.ByYear = make(map[string]*Bucket)
	var mind, maxd time.Time
	for _, msg := range a.MailBox.ById {
		payload.ArticleCount++
		if mind.IsZero() || msg.Date.Before(mind) {
			mind = msg.Date
		}
		if maxd.IsZero() || msg.Date.After(maxd) {
			maxd = msg.Date
		}
		year := msg.Date.Format("2006")
		bucket, ok := payload.ByYear[year]
		if !ok {
			bucket = &Bucket{Period: year}
			payload.ByYear[year] = bucket
		}
		bucket.Messages = append(bucket.Messages, msg)
	}
	payload.From = mind.Format("January 2, 2006")
	payload.Through = maxd.Format("January 2, 2006")

	return func(w http.ResponseWriter, r *http.Request) {
		a.render(w, r, payload, "layout", "index")
	}
}

func (a *App) handleMessage(w http.ResponseWriter, r *http.Request) {
	id := way.Param(r.Context(), "id")
	msg, ok := a.MailBox.ById[id]
	if !ok {
		if no, err := strconv.Atoi(id); err == nil {
			msg, ok = a.MailBox.ByLine[no]
		}
	}
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	a.render(w, r, msg, "layout", "message")
}

func (a *App) handleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := struct {
			Method string
			URL    string
		}{
			Method: r.Method,
			URL:    r.URL.Path,
		}
		a.render(w, r, payload, "layout", "not_found")
	}
}
