package mbox

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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

// spam, hide map[string]bool

func (h *Header) Parse() error {
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
	}

	return nil
}

// DateAsTime returns the "Date" header as time.Time.
// It understands the following formats:
func (h *Header) DateAsTime() (time.Time, error) {
	for _, layout := range []string{
		//	"Thu, 24 Mar 2011 20:09:09 -0700 (PDT)"
		"Mon, 2 Jan 2006 15:04:05 -0700 (MST)",
		// "09 Oct 2007 02:46:48 GMT"
		"02 Jan 2006 15:04:05 MST",
		// "10 May 2011 12:36:45 GMT"
		"2 Jan 2006 15:04:05 MST",
		// "02 Aug 2003 00:26:31 +0200"
		"02 Jan 2006 15:04:05 -0700",
		// "10 May 2011 08:53:15 -0400"
		"2 Jan 2006 15:04:05 -0700",
		// "Mon, 05 Jan 2009 13:44:02 -0600"
		"Mon, 02 Jan 2006 15:04:05 -0700",
		// "Tue, 4 Sep 2012 20:37:24 +0200"
		"Mon, 2 Jan 2006 15:04:05 -0700",
		// "Wed, 7 Jan 2009 07:26:42 GMT"
		"Mon, 02 Jan 2006 15:04:05 MST",
		// "Sun, 28 Dec 2008 22:43:09 GMT"
		"Mon, 2 Jan 2006 15:04:05 MST",
		// "15 Feb 01 17:44:09 GMT"
		"2 Jan 06 15:04:05 MST",
		// "Sat, 17 Feb 01 23:14:55"
		"Mon, 2 Jan 06 15:04:05",
		// "15 Dec 00 15:28:22 +0100"
		"2 Jan 06 15:04:05 -0700",
		// "2000/11/28"
		"2006/01/02",
		// "25 Mar 95 19:26:17"
		"2 Jan 06 15:04:05",
		// "Sat, 25 Mar 95 22:34:21 -0500"
		"Mon, 2 Jan 06 15:04:05 -0700",
		// "Mon, 27 Mar 1995 14:54:12"
		"Mon, 2 Jan 2006 15:04:05",
		// "Sat, 18 Mar 95 21:05:28 PDT"
		"Mon, 2 Jan 06 15:04:05 MST",
		// "Sun, 19 Mar 1995 08:37:28 LOCAL"
		"Mon, 2 Jan 2006 15:04:05 LOCAL",
		// "Sat, 11 Mar 1995 12:17:31 UNDEFINED"
		"Mon, 2 Jan 2006 15:04:05 UNDEFINED",
		// "5 Feb 1995 11:01 -0500"
		"2 Jan 2006 15:04 -0700",
		// "Thu, 12 Jan 1995 23:01"
		"Mon, 2 Jan 2006 15:05",
		// "17 Dec 1994 01:04 CST"
		"2 Jan 2006 15:04 MST",
		// "12 Oct 94 16:11:03 +"
		"2 Jan 06 15:04:05 +",
		// "Wed, 12 Oct 1994 09:35:51 Central"
		//"Mon, 2 Jan 2006 15:04:05 Central",
		// "Tue, 31 May 1994  13:46 MET"
		"Mon, 2 Jan 2006  15:04 MST",
		// "Mon, 23 May 94 15:51:27 -0700 (PDT)"
		"Mon, 2 Jan 06 15:04:05 -0700 (MST)",
		// "Thu, 02 Dec 93 19:50:54 est"
		"Mon, 02 Jan 06 15:04:05 MST",
		// "Monday, 12 Jul 1993 14:46:15 EDT"
		"Monday, 2 Jan 2006 15:04:05 MST",
	} {
		if t, err := time.Parse(layout, h.Date); err == nil {
			//if op := t.Format(layout); op != h.Date {
			//	if h.Date != strings.Replace(op, "+0000", "-0000", 1) {
			//		log.Printf("%d\n  fmt %q\nvalue %q\n   op %q\n", n+1, layout, h.Date, op)
			//	}
			//}
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("unknown layout %q", h.Date)
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
	s = strings.ReplaceAll(s, "<", " ")
	s = strings.ReplaceAll(s, ">", " ")
	s = strings.ReplaceAll(s, "}", " ")

	var ids []string
	for _, id := range strings.Fields(s) {
		if len(id) == 0 {
			continue
		}
		ids = append(ids, id)
	}

	return ids
}
