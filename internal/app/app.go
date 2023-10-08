// Package app implements an application server.
package app

import (
	"github.com/matryer/way"
	"github.com/mdhender/mbox/internal/stores/newsgroup"
	"sync"
)

type App struct {
	Host      string
	NewsGroup *newsgroup.NewsGroup
	NewSpam   struct {
		sync.Mutex
		AllowReports bool
		Posts        map[string]*newsgroup.Post
	}
	Port      string
	Router    *way.Router
	Templates string
}

func New(ng *newsgroup.NewsGroup, allowSpamReports bool) (*App, error) {
	a := &App{
		NewsGroup: ng,
		Port:      "8080",
		Router:    way.NewRouter(),
		Templates: "../templates",
	}
	a.NewSpam.AllowReports = allowSpamReports
	a.NewSpam.Posts = make(map[string]*newsgroup.Post)
	a.Router.HandleFunc("GET", "/", a.handleIndex())
	a.Router.HandleFunc("GET", "/corpus", a.handleCorpus())
	a.Router.HandleFunc("GET", "/corpus/:id", a.handleCorpusId)
	a.Router.HandleFunc("GET", "/from/:year", a.handleYear)
	a.Router.HandleFunc("GET", "/from/:year/:month", a.handleYearMonth)
	a.Router.HandleFunc("GET", "/from/:year/:month/:day", a.handleYearMonthDay)
	a.Router.HandleFunc("GET", "/posts", a.handlePostsSearch())
	a.Router.HandleFunc("GET", "/posts/:id", a.handlePosts)
	a.Router.HandleFunc("GET", "/spam", a.handleSpam)
	a.Router.NotFound = a.notFound()

	return a, nil
}
