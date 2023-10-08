package newsgroup

//// AddToCorpus add the words from a Chunk into the Corpus.
//func (ng *NewsGroup) AddToCorpus(p *Post, ch *chunk.Chunk) {
//	words := make(map[string]int)
//	ng.Corpus.Documents[p.ShaId] = words
//	for _, line := range ch.Body {
//		// split on any non-letter/non-number rune.
//		for _, word := range bytes.FieldsFunc(line, func(r rune) bool {
//			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
//		}) {
//			word := strings.ToLower(string(word))
//			if _, ok := ng.Corpus.StopWords[word]; ok {
//				// filter out the stop word
//				continue
//			}
//			words[word] = words[word] + 1
//		}
//	}
//
//	// add the words to the index
//	for word := range words {
//		ng.Corpus.Index[word] = append(ng.Corpus.Index[word], p)
//	}
//}
