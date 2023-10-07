package mbox

//func (mb *MailBox) MakeCorpus() {
//	// words and stuff
//	for _, msg := range mb.Messages {
//		msg.ParseWords()
//	}
//
//	// corpus
//	for _, msg := range mb.Messages {
//		for k, v := range msg.Words {
//			if isword(k) {
//				mb.Corpus[k] = mb.Corpus[k] + v
//			}
//		}
//	}
//}
//
//func isword(s string) bool {
//	for _, r := range s {
//		if !unicode.IsLetter(r) {
//			return false
//		}
//	}
//	return true
//}
