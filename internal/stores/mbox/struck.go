package mbox

func (mb *MailBox) FlagStruck() {
	// code to display header for struck messages. should be commented out.
	for _, msg := range mb.Messages {
		if msg.Struck {
			continue
		}
	}
}

var struck = map[string]bool{}
