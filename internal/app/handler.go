package app

import (
	"github.com/matryer/way"
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
	var payload struct {
		ArticleCount int
		From         string
		Through      string
		Years        map[string]int
	}
	payload.Years = a.NewsGroup.Posts.Years
	var mind, maxd time.Time
	for _, post := range a.NewsGroup.Posts.ByShaId {
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
	}
	payload.From = mind.Format("January 2, 2006")
	payload.Through = maxd.Format("January 2, 2006")

	return func(w http.ResponseWriter, r *http.Request) {
		a.render(w, r, payload, "layout", "index")
	}
}

func (a *App) handlePost(w http.ResponseWriter, r *http.Request) {
	id := way.Param(r.Context(), "id")
	post, ok := a.NewsGroup.Posts.ByShaId[id]
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

func (a *App) handleNotFound(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Method string
		URL    string
	}{
		Method: r.Method,
		URL:    r.URL.Path,
	}
	a.render(w, r, payload, "layout", "not_found")
}

func (a *App) handleYear(w http.ResponseWriter, r *http.Request) {
	year := way.Param(r.Context(), "year")
	bucket, ok := a.NewsGroup.Posts.ByPeriod[year]
	if !ok {
		log.Printf("[app] year %q not found\n", year)
		a.handleNotFound(w, r)
		return
	}
	a.render(w, r, bucket, "layout", "from_period")
}

func (a *App) handleYearMonth(w http.ResponseWriter, r *http.Request) {
	year := way.Param(r.Context(), "year")
	month := way.Param(r.Context(), "month")
	bucket, ok := a.NewsGroup.Posts.ByPeriod[year+"/"+month]
	if !ok {
		log.Printf("[app] year %q month %q not found\n", year, month)
		a.handleNotFound(w, r)
		return
	}
	a.render(w, r, bucket, "layout", "from_period")
}

func (a *App) handleYearMonthDay(w http.ResponseWriter, r *http.Request) {
	year := way.Param(r.Context(), "year")
	month := way.Param(r.Context(), "month")
	day := way.Param(r.Context(), "day")
	bucket, ok := a.NewsGroup.Posts.ByPeriod[year+"/"+month+"/"+day]
	if !ok {
		log.Printf("[app] year %q month %q day %q not found\n", year, month, day)
		a.handleNotFound(w, r)
		return
	}
	a.render(w, r, bucket, "layout", "from_period")
}

func (a *App) notFound() http.HandlerFunc {
	return a.handleNotFound
}
