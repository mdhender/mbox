package app

import (
	"fmt"
	"github.com/matryer/way"
	"github.com/mdhender/mbox/internal/stores/newsgroup"
	"log"
	"net/http"
	"strings"
	"time"
)

func (a *App) handleCorpus() http.HandlerFunc {
	payload := struct {
		Posts map[string]*newsgroup.Post
		Index map[string][]*newsgroup.Post
	}{
		Posts: a.NewsGroup.Posts.ByShaId,
		Index: a.NewsGroup.Corpus.Index,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[corpus] dump %d\n", len(payload.Posts))
		a.render(w, r, payload, "layout", "corpus")
	}
}

func (a *App) handleCorpusId(w http.ResponseWriter, r *http.Request) {
	id := way.Param(r.Context(), "id")
	post, ok := a.NewsGroup.Posts.ByShaId[id]
	if !ok {
		log.Printf("[app] corpus post %q not found\n", id)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	a.render(w, r, post, "layout", "corpus_id")
}

func (a *App) handleIndex() http.HandlerFunc {
	var payload struct {
		AllowSpamReporting bool
		ArticleCount       int
		From               string
		Search             string
		Through            string
		Years              map[string]int
	}
	payload.AllowSpamReporting = a.NewSpam.AllowReports
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

// post may be a simple index or a complicated query
func (a *App) handlePosts(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		AllowSpamReporting bool
		Post               *newsgroup.Post
	}
	payload.AllowSpamReporting = a.NewSpam.AllowReports

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
	payload.Post = post

	if pSpamFlag := r.URL.Query().Get("spam"); a.NewSpam.AllowReports && pSpamFlag != "" {
		a.NewSpam.Lock()
		defer a.NewSpam.Unlock()
		switch pSpamFlag {
		case "false":
			if _, ok := a.NewSpam.Posts[post.ShaId]; ok {
				post.Spam = false
				post.Subject = strings.TrimPrefix(post.Subject, "(spam) ")
				delete(a.NewSpam.Posts, post.ShaId)
			}
		case "true":
			if post.Spam {
				// already tagged, so do nothing
			} else if _, ok := a.NewSpam.Posts[post.ShaId]; ok {
				// already in the new spam group?
			} else {
				post.Spam = true
				if !strings.HasPrefix(post.Subject, "(spam) ") {
					post.Subject = "(spam) " + post.Subject
				}
				a.NewSpam.Posts[post.ShaId] = post
			}
		}
		http.Redirect(w, r, fmt.Sprintf("/posts/%s", post.ShaId), http.StatusSeeOther)
		return
	}
	a.render(w, r, payload, "layout", "post")
}

func (a *App) handlePostsSearch() http.HandlerFunc {
	type payload struct {
		AllowSpamReporting bool
		Search             string
		Posts              []*newsgroup.Post
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var p payload
		p.AllowSpamReporting = a.NewSpam.AllowReports
		p.Search = r.URL.Query().Get("q")
		// log.Printf("[posts] pSearch %+v\n", p.Search)

		if p.Search != "" {
			// "compliment blessed" firestorm Morghoul perceval dean
			for _, post := range a.NewsGroup.SearchPosts(p.Search) { //
				p.Posts = append(p.Posts, post)
				// log.Printf("[search] post http://localhost:8080/posts/%s\n", post.ShaId)
			}
		}

		a.render(w, r, p, "layout", "posts_search")
	}
}

func (a *App) handleSpam(w http.ResponseWriter, r *http.Request) {
	a.NewSpam.Lock()
	defer a.NewSpam.Unlock()

	var payload struct {
		Lines int    // number of spams
		Body  string // list of spam
	}
	payload.Lines = len(a.NewSpam.Posts) + 2
	for _, post := range a.NewSpam.Posts {
		payload.Body += fmt.Sprintf("\n%q: true,", post.Id)
	}
	payload.Body += "\n"
	a.render(w, r, payload, "layout", "spam")
}

func (a *App) handleYear(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		AllowSpamReporting bool
		Bucket             *newsgroup.Bucket
	}
	payload.AllowSpamReporting = a.NewSpam.AllowReports

	year := way.Param(r.Context(), "year")
	bucket, ok := a.NewsGroup.Posts.ByPeriod[year]
	if !ok {
		log.Printf("[app] year %q not found\n", year)
		a.handleNotFound(w, r)
		return
	}
	payload.Bucket = bucket
	a.render(w, r, payload, "layout", "from_period")
}

func (a *App) handleYearMonth(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		AllowSpamReporting bool
		Bucket             *newsgroup.Bucket
	}
	payload.AllowSpamReporting = a.NewSpam.AllowReports

	year := way.Param(r.Context(), "year")
	month := way.Param(r.Context(), "month")
	bucket, ok := a.NewsGroup.Posts.ByPeriod[year+"/"+month]
	if !ok {
		log.Printf("[app] year %q month %q not found\n", year, month)
		a.handleNotFound(w, r)
		return
	}
	payload.Bucket = bucket
	a.render(w, r, payload, "layout", "from_period")
}

func (a *App) handleYearMonthDay(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		AllowSpamReporting bool
		Bucket             *newsgroup.Bucket
	}
	payload.AllowSpamReporting = a.NewSpam.AllowReports

	year := way.Param(r.Context(), "year")
	month := way.Param(r.Context(), "month")
	day := way.Param(r.Context(), "day")
	bucket, ok := a.NewsGroup.Posts.ByPeriod[year+"/"+month+"/"+day]
	if !ok {
		log.Printf("[app] year %q month %q day %q not found\n", year, month, day)
		a.handleNotFound(w, r)
		return
	}
	payload.Bucket = bucket
	a.render(w, r, payload, "layout", "from_period")
}

func (a *App) notFound() http.HandlerFunc {
	return a.handleNotFound
}
