package mbox

import (
	"regexp"
)

func (mb *MailBox) Parse(lines [][]byte) {
	// every message starts with a blank line followed by a From with ID.
	som := regexp.MustCompile("^From -?[0-9]+$")

	// extract all the messages in the file
	var msg *Message
	for n, line := range lines {
		if som.Find(line) != nil {
			msg = &Message{
				Start: n,
			}
			mb.Messages = append(mb.Messages, msg)
		}

		if msg != nil {
			msg.Raw = append(msg.Raw, line)
		}
	}

	// split the message header and body
	for _, msg := range mb.Messages {
		msg.ExtractHeaderAndBody()
	}
}

func collectHeaderKeys(msgs []*Message) (map[string]int, error) {
	headers := make(map[string]int)
	for _, msg := range msgs {
		keys, err := msg.Header.Keys()
		if err != nil {
			return nil, err
		}
		for _, key := range keys {
			headers[key] = headers[key] + 1
		}
	}
	return headers, nil
}
