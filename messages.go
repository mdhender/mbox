package main

import (
	"bytes"
	"os"
	"regexp"
)

var mbox = map[string]*Message{}

func read(path string) ([][]byte, error) {
	data, err := os.ReadFile("rec.games.pbm.mbox")
	if err != nil {
		return nil, err
	}
	// prepend a blank line to make splitting messages simpler
	data = append([]byte{'\n'}, data...)
	// split into lines and trim any carriage-returns
	lines := bytes.Split(data, []byte{'\n'})
	for i := 0; i < len(lines); i++ {
		lines[i] = bytes.TrimRight(lines[i], "\r")
	}
	return lines, nil
}

func split(lines [][]byte) ([]*Message, error) {
	// every message starts with a blank line followed by a From with ID.
	som := regexp.MustCompile("^From -?[0-9]+$")

	//state := "looking for message"
	var msgs []*Message
	var msg *Message

	// process all the messages in the file
	for n, line := range lines {
		if som.Find(line) != nil {
			msg = &Message{
				Start: n,
			}
			msgs = append(msgs, msg)
		}

		if msg != nil {
			msg.Raw = append(msg.Raw, line)
		}
	}

	return msgs, nil
}

func CollectHeaderKeys(msgs []*Message) (map[string]int, error) {
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
