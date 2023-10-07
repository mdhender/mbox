package memory

import "time"

type Message struct {
	Id           string
	Spam         bool
	Struck       bool
	Subject      string
	Date         time.Time
	Body         string
	References   []*Message
	ReferencedBy []*Message
}
