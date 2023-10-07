package app

import (
	"github.com/matryer/way"
	"github.com/mdhender/mbox/internal/stores/newsgroup"
	"log"
	"net/http"
	"time"
)

func (a *App) handleCorpus(w http.ResponseWriter, r *http.Request) {
	log.Printf("[corpus] dump %d\n", len(a.NewsGroup.Corpus.Words))
	payload := struct {
		Corpus map[string]int
	}{
		Corpus: a.NewsGroup.Corpus.Words,
	}
	a.render(w, r, payload, "layout", "corpus")
}

func (a *App) handleIndex() http.HandlerFunc {
	type Bucket struct {
		Period string
		Posts  []*newsgroup.Post
	}
	var payload struct {
		ArticleCount int
		From         string
		Through      string
		ByYear       map[string]*Bucket
	}
	payload.ByYear = make(map[string]*Bucket)
	var mind, maxd time.Time
	for _, post := range a.NewsGroup.Posts.ById {
		// don't include missing posts
		if post.Missing {
			continue
		}
		payload.ArticleCount++
		if mind.IsZero() || post.Date.Before(mind) {
			mind = post.Date
		}
		if maxd.IsZero() || post.Date.After(maxd) {
			maxd = post.Date
		}
		year := post.Date.Format("2006")
		if year == "0001" {
			log.Printf("hey post %q is year %s\n", post.Id, year)
		}
		bucket, ok := payload.ByYear[year]
		if !ok {
			bucket = &Bucket{Period: year}
			payload.ByYear[year] = bucket
		}
		bucket.Posts = append(bucket.Posts, post)
	}
	payload.From = mind.Format("January 2, 2006")
	payload.Through = maxd.Format("January 2, 2006")

	return func(w http.ResponseWriter, r *http.Request) {
		a.render(w, r, payload, "layout", "index")
	}
}

func (a *App) handlePost(w http.ResponseWriter, r *http.Request) {
	id := way.Param(r.Context(), "id")
	post, ok := a.NewsGroup.Posts.ById[id]
	if ok {
		log.Printf("[app] found post %q by id %q\n", post.Id, id)
	}
	if !ok {
		post, ok = a.NewsGroup.Posts.ByLineNo[id]
		if ok {
			log.Printf("[app] found post %q by line number %q\n", post.Id, id)
		}
	}
	if !ok {
		log.Printf("[app] post %q not found\n", id)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	a.render(w, r, post, "layout", "post")
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
