package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Header struct {
	Id     string
	Date   string
	From   string
	Google struct {
		Thread     string
		NewGroupId bool
	}
	Lines      int
	References struct {
		Ids      []string //
		Messages []*Message
		Unknown  []string
	}
	Subject string
	Text    [][]byte
}

func (h *Header) Keys() ([]string, error) {
	var keys []string
	for _, kvp := range h.Text {
		key, value, found := strings.Cut(string(kvp), ":")
		if found {
			// may be a bug to force the keys to lower-case, but meh
			keys = append(keys, strings.ToLower(key))
		} else if value != "" {
			panic("assert(value != '')")
		}
	}
	return keys, nil
}

func (h *Header) Parse(spam, hide map[string]bool) error {
	for _, text := range h.Text {
		key, value, found := strings.Cut(string(text), ":")
		if !found {
			return fmt.Errorf("invalid header line")
		}
		key, value = strings.ToLower(key), strings.TrimSpace(value)
		switch key {
		case "date":
			h.Date = value
		case "from":
			h.From = value
		case "lines":
			lines, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("invalid header lines: %w", err)
			}
			h.Lines = lines
		case "message-id":
			if strings.HasPrefix(value, "<") && strings.HasSuffix(value, ">") {
				h.Id = value[1 : len(value)-1]
			} else if strings.HasPrefix(value, "<") && strings.HasSuffix(value, ">#1/1") {
				h.Id = value[1 : len(value)-5]
			} else {
				return fmt.Errorf("invalid message-id %q", value)
			}
		case "references":
			h.References.Ids = parseReferences(value)
		case "subject":
			h.Subject = value
		case "x-google-newgroupid":
			switch value {
			case "yes":
				h.Google.NewGroupId = true
			default:
				return fmt.Errorf("unknown x-google-newgroupid %q", value)
			}
		case "x-google-thread":
			h.Google.Thread = value
		default:
			// ignore key
		}
	}

	if h.Id == "" {
		return fmt.Errorf("missing mesage-id")
	} else if mbox[h.Id] != nil {
		return fmt.Errorf("duplicate message id %q", h.Id)
	}

	if spam[h.Id] {
		h.Subject = "Message has been flagged as spam."
	} else if hide[h.Id] {
		h.Subject = "Message has been taken down due to request."
	}

	return nil
}

func (h *Header) LinkReferences(mbox map[string]*Message) error {
	for _, id := range h.References.Ids {
		if msg, ok := mbox[id]; !ok {
			h.References.Unknown = append(h.References.Unknown, id)
		} else {
			h.References.Messages = append(h.References.Messages, msg)
		}
	}
	return nil
}

// references should be a list of message id values separated by spaces
func parseReferences(s string) []string {
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.TrimSpace(s)

	var ids []string
	for _, id := range strings.Split(s, " ") {
		if len(id) < 4 {
			continue
		}
		if id[0] == '<' {
			id = id[1:]
		}
		if id[len(id)-1] == '>' || id[len(id)-1] == '}' {
			id = id[:len(id)-1]
		}
		ids = append(ids, id)
	}
	return ids
}
