package newsgroup

import (
	"fmt"
	"github.com/mdhender/mbox/internal/chunk"
	"log"
)

func (ng *NewsGroup) Parse(ch *chunk.Chunk) error {
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
		return fmt.Errorf("post %q: %w", string(ch.From[5:]), err)
	}

	if p.Id == "" {
		log.Printf("[post] %q: missing id", string(ch.From[5:]))
		p.Id = string(ch.From[5:])
	}

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
		return fmt.Errorf("post %q: %w", string(ch.From[5:]), err)
	}

	if ng.Posts.ById[p.Id] != nil {
		return fmt.Errorf("post %q: duplicate id %q", string(ch.From[5:]), p.Id)
	}
	ng.Posts.ById[p.Id] = p
	ng.Posts.ByLineNo[fmt.Sprintf("%d", p.LineNo)] = p

	return nil
}
