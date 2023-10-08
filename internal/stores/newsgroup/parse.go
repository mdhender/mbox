package newsgroup

import (
	"fmt"
	"github.com/mdhender/mbox/internal/chunk"
	"log"
)

func (ng *NewsGroup) Parse(ch *chunk.Chunk, createCorpus bool) (*Post, error) {
	p := &Post{
		Keys:         make(map[string][]string),
		LineNo:       ch.Line,
		References:   make(map[string]*Post),
		ReferencedBy: make(map[string]*Post),
		Sender:       "(missing sender)",
		Subject:      "(Missing Subject Line)",
	}

	// parse the header
	err := p.ParseHeader(ch)
	if err != nil {
		return nil, fmt.Errorf("post %q: %w", string(ch.From[5:]), err)
	}

	if p.Id == "" {
		log.Printf("[post] %q: missing id", string(ch.From[5:]))
		p.Id = string(ch.From[5:])
	}
	p.ShaId = sha1sum(p.Id)

	// flag spam and stuck messages
	if p.Spam = ng.Posts.Spam[p.Id]; p.Spam {
		p.Subject = "Message has been flagged as spam."
	}
	if p.Struck = ng.Posts.Struck[p.Id]; p.Struck {
		p.Subject = "Message has been taken down due to request."
	}

	// parse the body
	err = p.ParseBody(ch)
	if err != nil {
		return nil, fmt.Errorf("post %q: %w", string(ch.From[5:]), err)
	}
	// don't index if spam or struck
	if p.Spam || p.Struck {
		return p, nil
	}

	if ng.Posts.ById[p.Id] != nil {
		return nil, fmt.Errorf("post %q: duplicate id %q", string(ch.From[5:]), p.Id)
	}
	ng.Posts.ById[p.Id] = p
	ng.Posts.ByLineNo[fmt.Sprintf("%d", p.LineNo)] = p
	ng.Posts.ByShaId[p.ShaId] = p

	// add this post to all the buckets
	year := p.Date.Format("2006")
	ng.Posts.Years[year] = ng.Posts.Years[year] + 1
	yearBucket, ok := ng.Posts.ByPeriod[year]
	if !ok {
		yearBucket = &Bucket{
			Up:         "/",
			Period:     year,
			SubPeriods: make(map[string]*Bucket),
			Posts:      []*Post{p},
		}
		ng.Posts.ByPeriod[year] = yearBucket
	}
	yearMonth := p.Date.Format("2006/01")
	monthBucket, ok := ng.Posts.ByPeriod[yearMonth]
	if !ok {
		monthBucket = &Bucket{
			Up:         "/from/" + year,
			Period:     yearMonth,
			SubPeriods: make(map[string]*Bucket),
			Posts:      []*Post{p},
		}
		ng.Posts.ByPeriod[yearMonth] = monthBucket
		yearBucket.SubPeriods[monthBucket.Period] = monthBucket
	}
	yearMonthDay := p.Date.Format("2006/01/02")
	dayBucket, ok := ng.Posts.ByPeriod[yearMonthDay]
	if !ok {
		dayBucket = &Bucket{
			Up:         "/from/" + yearMonth,
			Period:     yearMonthDay,
			SubPeriods: make(map[string]*Bucket),
			Posts:      []*Post{p},
		}
		ng.Posts.ByPeriod[yearMonthDay] = dayBucket
		monthBucket.SubPeriods[dayBucket.Period] = dayBucket
	}
	dayBucket.Posts = append(dayBucket.Posts, p)
	p.Up = "/from/" + yearMonth

	if createCorpus {
		p.Words = ch.Words(ng.Corpus.StopWords)
		for word := range p.Words {
			ng.Corpus.Index[word] = append(ng.Corpus.Index[word], p)
		}
	}

	return p, nil
}
