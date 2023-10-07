package main

import (
	"github.com/matryer/way"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type App struct {
	messages struct {
		byId   map[string]*Message
		byLine map[int]*Message
		corpus map[string]int
	}
	router    *way.Router
	templates string
}

func (a *App) render(w http.ResponseWriter, r *http.Request, data any, names ...string) {
	var files []string
	for _, name := range names {
		files = append(files, filepath.Join(a.templates, name+".gohtml"))
	}

	t, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("%s %s: render: parse: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("%s %s: render: execute: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (a *App) handleCorpus(w http.ResponseWriter, r *http.Request) {
	log.Printf("[corpus] dump %d\n", len(a.messages.corpus))
	payload := struct {
		Corpus map[string]int
	}{
		Corpus: a.messages.corpus,
	}
	a.render(w, r, payload, "layout", "corpus")
}

func (a *App) handleIndex() http.HandlerFunc {
	type Bucket struct {
		Period   string
		Messages []*Message
	}
	var payload struct {
		ArticleCount int
		From         string
		Through      string
		ByYear       map[string]*Bucket
	}
	payload.ByYear = make(map[string]*Bucket)
	var mind, maxd time.Time
	for _, msg := range a.messages.byId {
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
	msg, ok := a.messages.byId[id]
	if !ok {
		if no, err := strconv.Atoi(id); err == nil {
			msg, ok = a.messages.byLine[no]
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
