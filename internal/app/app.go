// Package app implements an application server.
package app

import (
	"github.com/matryer/way"
	"github.com/mdhender/mbox/internal/stores/mbox"
)

type App struct {
	MailBox   *mbox.MailBox
	Router    *way.Router
	Templates string
}

func New(mb *mbox.MailBox) (*App, error) {
	a := &App{
		Router:    way.NewRouter(),
		Templates: "../templates",
	}
	a.MailBox = mb
	a.Router.HandleFunc("GET", "/", a.handleIndex())
	a.Router.HandleFunc("GET", "/corpus", a.handleCorpus)
	a.Router.HandleFunc("GET", "/messages/:id", a.handleMessage)
	a.Router.NotFound = a.handleNotFound()

	return a, nil
}
