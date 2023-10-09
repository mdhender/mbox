package app

import (
	"github.com/matryer/way"
	"log"
	"net/http"
	"sort"
	"time"
)

type Index struct {
	ArticleCount int
	From         string
	Through      string
	Years        []*Period
}

type Post struct {
	Id           string
	Url          string
	Spam         bool
	Struck       bool
	From         string
	Subject      string
	Date         string
	Lines        int
	Body         string
	References   []Reference // list of id
	ReferencedBy []Reference // list of id
	Parent       string      // url of parent post
}

type PostsCollection struct {
	Name   string
	Parent string
	Posts  []*Post
}

type Reference struct {
	Url     string
	From    string
	Subject string
	Date    string
}

type Period struct {
	Name  string
	Url   string
	Count int
}

type Bucket struct {
	Name     string
	Url      string
	Parent   string
	Count    int
	Children []*Bucket
}

func (a *App) handleIndex() http.HandlerFunc {
	var payload Index
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
	for year, count := range a.NewsGroup.Posts.Years {
		payload.Years = append(payload.Years, &Period{
			Name:  year,
			Count: count,
			Url:   "/from/" + year,
		})
	}
	sort.Slice(payload.Years, func(i, j int) bool {
		return payload.Years[i].Name < payload.Years[j].Name
	})

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
	var payload Post

	id := way.Param(r.Context(), "id")
	post, ok := a.NewsGroup.Posts.ByShaId[id]
	if !ok {
		log.Printf("[app] post %q not found\n", id)
		a.handleNotFound(w, r)
		return
	}
	log.Printf("[app] found post %q by id %q\n", post.Id, id)
	payload = Post{
		Id:      post.ShaId,
		Url:     "/posts/" + post.ShaId,
		Spam:    post.Spam,
		Struck:  post.Struck,
		From:    post.Sender,
		Subject: post.Subject,
		Date:    post.Date.Format(time.RFC1123Z),
		Lines:   post.Lines,
		Body:    post.Body,
		Parent:  post.Date.Format("/from/2006/01/02"),
	}
	for _, ref := range post.References {
		if ref.Subject != "** missing post **" {
			payload.References = append(payload.References, Reference{
				Url:     "/posts/" + ref.ShaId,
				From:    ref.Sender,
				Subject: ref.Subject,
				Date:    ref.Date.Format(time.RFC1123Z),
			})
		}
	}
	for _, ref := range post.ReferencedBy {
		if ref.Subject != "** missing post **" {
			payload.ReferencedBy = append(payload.ReferencedBy, Reference{
				Url:     "/posts/" + ref.ShaId,
				From:    ref.Sender,
				Subject: ref.Subject,
				Date:    ref.Date.Format(time.RFC1123Z),
			})
		}
	}

	a.render(w, r, payload, "layout", "post")
}

func (a *App) handleYear(w http.ResponseWriter, r *http.Request) {
	year := way.Param(r.Context(), "year")
	payload := Bucket{Name: year, Parent: "/posts"}
	bucket, ok := a.NewsGroup.Posts.ByPeriod[year]
	if !ok {
		log.Printf("[app] year %q not found\n", year)
		a.handleNotFound(w, r)
		return
	}
	for _, child := range bucket.SubPeriods {
		payload.Children = append(payload.Children, &Bucket{
			Name:  child.Period,
			Url:   "/from/" + child.Period,
			Count: child.Count(),
		})
	}
	sort.Slice(payload.Children, func(i, j int) bool {
		return payload.Children[i].Name < payload.Children[j].Name
	})
	a.render(w, r, payload, "layout", "from_yyyy")
}

func (a *App) handleYearMonth(w http.ResponseWriter, r *http.Request) {
	year := way.Param(r.Context(), "year")
	month := way.Param(r.Context(), "month")
	payload := Bucket{Name: year + "/" + month, Parent: "/from/" + year}
	bucket, ok := a.NewsGroup.Posts.ByPeriod[payload.Name]
	if !ok {
		log.Printf("[app] year %q month %q not found\n", year, month)
		a.handleNotFound(w, r)
		return
	}
	for _, child := range bucket.SubPeriods {
		payload.Children = append(payload.Children, &Bucket{
			Name:  child.Period,
			Url:   "/from/" + child.Period,
			Count: child.Count(),
		})
	}
	sort.Slice(payload.Children, func(i, j int) bool {
		return payload.Children[i].Name < payload.Children[j].Name
	})
	a.render(w, r, payload, "layout", "from_yyyy_mm")
}

func (a *App) handleYearMonthDay(w http.ResponseWriter, r *http.Request) {
	year := way.Param(r.Context(), "year")
	month := way.Param(r.Context(), "month")
	day := way.Param(r.Context(), "day")
	payload := PostsCollection{
		Name:   year + "/" + month + "/" + day,
		Parent: "/from/" + year + "/" + month,
	}
	bucket, ok := a.NewsGroup.Posts.ByPeriod[payload.Name]
	if !ok {
		log.Printf("[app] year %q month %q day %q not found\n", year, month, day)
		a.handleNotFound(w, r)
		return
	}
	for _, post := range bucket.Posts {
		payload.Posts = append(payload.Posts, &Post{
			Url:     "/posts/" + post.ShaId,
			From:    post.Sender,
			Subject: post.Subject,
			Date:    post.Date.Format("15:04:05"),
		})
	}
	a.render(w, r, payload, "layout", "from_yyyy_mm_dd")
}

func (a *App) notFound() http.HandlerFunc {
	return a.handleNotFound
}
