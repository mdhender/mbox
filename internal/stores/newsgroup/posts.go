package newsgroup

import (
	"fmt"
	"github.com/mdhender/mbox/internal/chunk"
	"log"
	"strings"
	"time"
)

// Post is a single posting to the newsgroup
type Post struct {
	Id           string              // unique ID from the "From " block header
	ShaId        string              // SHA-1 hash of the Id
	Body         string              // body of the posting
	Date         time.Time           // time post was added to the newsgroup
	Error        error               // any error parsing the message
	Keys         map[string][]string // unknown (or ignored) keys and values
	Lines        int                 // number of lines in post body
	LineNo       int                 // line number from original mbox file
	Missing      bool                // true if the original message is missing
	References   map[string]*Post    // posts this post references
	ReferencedBy map[string]*Post    // posts referring to this post
	Sender       string              // e-mail address of person sending the post
	Spam         bool                // post is considered spam
	Struck       bool                // post is struck for copyright or ownership
	Subject      string              // subject of post
}

// ParseBody populates body from the input Chunk.
// Assumes the Spam and Struck flags have been set in the header.
func (p *Post) ParseBody(ch *chunk.Chunk) error {
	if p.Spam || p.Struck {
		p.Lines, p.Body = 3, "Message content has been removed."
		return nil
	} else if p.Missing {
		p.Subject = "(missing post)"
		p.Lines, p.Body = 3, "Unable to find original posting.\n"
		return nil
	}

	sb := strings.Builder{}
	for _, line := range ch.Body {
		sb.Write(line)
		sb.WriteByte('\n')
	}
	p.Lines, p.Body = len(ch.Body)+2, sb.String()
	return nil
}

// ParseHeader updates header values from the input Chunk.
func (p *Post) ParseHeader(ch *chunk.Chunk) error {
	debug := false // bytes.Equal(ch.From, []byte("From -3534941848242442294"))
	for _, text := range ch.Header {
		key, value, found := strings.Cut(string(text), ":")
		if !found {
			return fmt.Errorf("invalid header line")
		}
		key, value = strings.ToLower(key), strings.TrimSpace(value)
		if debug {
			log.Printf("[parse header] %q key %q value %q\n", string(ch.From), key, value)
		}
		switch key {
		case "date":
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
				if t, err := time.Parse(layout, value); err == nil {
					//if op := t.Format(layout); op != h.Date {
					//	if h.Date != strings.Replace(op, "+0000", "-0000", 1) {
					//		log.Printf("%d\n  fmt %q\nvalue %q\n   op %q\n", n+1, layout, h.Date, op)
					//	}
					//}
					p.Date = t.UTC()
					break
				}
			}
			if p.Date.IsZero() {
				return fmt.Errorf("unknown layout %q", value)
			}
		case "from":
			p.Sender = value
		case "message-id":
			if strings.HasPrefix(value, "<") && strings.HasSuffix(value, ">") {
				p.Id = value[1 : len(value)-1]
			} else if strings.HasPrefix(value, "<") && strings.HasSuffix(value, ">#1/1") {
				p.Id = value[1 : len(value)-5]
			} else {
				return fmt.Errorf("invalid message-id %q", value)
			}
		case "references":
			for _, id := range strings.Fields(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(value, "<", " "), ">", " "), "}", " ")) {
				if len(id) == 0 {
					continue
				}
				p.References[id] = nil
			}
		case "subject":
			p.Subject = value
		default:
			p.Keys[key] = append(p.Keys[key], value)
		}
	}

	return nil
}
