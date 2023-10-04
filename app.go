package main

import (
	"github.com/matryer/way"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

type App struct {
	messages struct {
		byId   map[string]*Message
		byLine map[int]*Message
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
