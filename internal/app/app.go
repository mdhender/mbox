// Package app implements an application server.
package app

import (
	"github.com/matryer/way"
	"github.com/mdhender/mbox/internal/stores/newsgroup"
)

type App struct {
	Host      string
	NewsGroup *newsgroup.NewsGroup
	Port      string
	Router    *way.Router
	Templates string
}

func New(ng *newsgroup.NewsGroup) (*App, error) {
	a := &App{
		NewsGroup: ng,
		Port:      "8080",
		Router:    way.NewRouter(),
		Templates: "../templates",
	}
	a.Router.HandleFunc("GET", "/", a.handleIndex())
	a.Router.HandleFunc("GET", "/corpus", a.handleCorpus)
	a.Router.HandleFunc("GET", "/posts/:id", a.handlePost)
	a.Router.NotFound = a.handleNotFound()

	return a, nil
}
