package newsgroup

import (
	"crypto/sha1"
	"encoding/base64"
	"log"
)

type NewsGroup struct {
	Corpus struct {
		// Words is a frequency count of words in the corpus.
		Words map[string]int
		// StopWords is a list of common usenet words to exclude from the indexing.
		StopWords map[string]bool
	}
	Posts struct {
		ById     map[string]*Post
		ByLineNo map[string]*Post
		ByPeriod map[string]*Bucket
		ByShaId  map[string]*Post
		Spam     map[string]bool
		Struck   map[string]bool
		Years    map[string]int
	}
}

type Bucket struct {
	Up         string // url of parent period
	Period     string
	SubPeriods map[string]*Bucket
	Posts      []*Post
}

func New() *NewsGroup {
	ng := &NewsGroup{}
	ng.Corpus.Words = make(map[string]int)
	ng.Corpus.StopWords = map[string]bool{
		"a":          true,
		"about":      true,
		"above":      true,
		"after":      true,
		"again":      true,
		"against":    true,
		"all":        true,
		"am":         true,
		"an":         true,
		"and":        true,
		"any":        true,
		"are":        true,
		"aren't":     true,
		"as":         true,
		"at":         true,
		"be":         true,
		"because":    true,
		"been":       true,
		"before":     true,
		"being":      true,
		"below":      true,
		"between":    true,
		"both":       true,
		"but":        true,
		"by":         true,
		"can":        true,
		"can't":      true,
		"cannot":     true,
		"could":      true,
		"couldn't":   true,
		"did":        true,
		"didn't":     true,
		"do":         true,
		"does":       true,
		"doesn't":    true,
		"doing":      true,
		"don't":      true,
		"down":       true,
		"during":     true,
		"each":       true,
		"ery":        true,
		"few":        true,
		"for":        true,
		"from":       true,
		"further":    true,
		"had":        true,
		"hadn't":     true,
		"has":        true,
		"hasn't":     true,
		"have":       true,
		"haven't":    true,
		"having":     true,
		"he":         true,
		"he'd":       true,
		"he'll":      true,
		"he's":       true,
		"her":        true,
		"here":       true,
		"here's":     true,
		"hers":       true,
		"herself":    true,
		"him":        true,
		"himself":    true,
		"his":        true,
		"how":        true,
		"how's":      true,
		"i":          true,
		"i'd":        true,
		"i'll":       true,
		"i'm":        true,
		"i've":       true,
		"if":         true,
		"in":         true,
		"into":       true,
		"is":         true,
		"isn't":      true,
		"it":         true,
		"it's":       true,
		"its":        true,
		"itself":     true,
		"let's":      true,
		"me":         true,
		"more":       true,
		"most":       true,
		"mustn't":    true,
		"my":         true,
		"myself":     true,
		"no":         true,
		"nor":        true,
		"not":        true,
		"of":         true,
		"off":        true,
		"on":         true,
		"once":       true,
		"only":       true,
		"or":         true,
		"other":      true,
		"ought":      true,
		"our":        true,
		"ours":       true,
		"ourselves":  true,
		"out":        true,
		"over":       true,
		"own":        true,
		"said":       true,
		"same":       true,
		"say":        true,
		"says":       true,
		"shall":      true,
		"shan't":     true,
		"she":        true,
		"she'd":      true,
		"she'll":     true,
		"she's":      true,
		"should":     true,
		"shouldn't":  true,
		"so":         true,
		"some":       true,
		"such":       true,
		"than":       true,
		"that":       true,
		"that's":     true,
		"the":        true,
		"their":      true,
		"theirs":     true,
		"them":       true,
		"themselves": true,
		"then":       true,
		"there":      true,
		"there's":    true,
		"these":      true,
		"they":       true,
		"they'd":     true,
		"they'll":    true,
		"they're":    true,
		"they've":    true,
		"this":       true,
		"those":      true,
		"through":    true,
		"to":         true,
		"too":        true,
		"under":      true,
		"until":      true,
		"up":         true,
		"upon":       true,
		"us":         true,
		"was":        true,
		"wasn't":     true,
		"we":         true,
		"we'd":       true,
		"we'll":      true,
		"we're":      true,
		"we've":      true,
		"were":       true,
		"weren't":    true,
		"what":       true,
		"what's":     true,
		"when":       true,
		"when's":     true,
		"where":      true,
		"where's":    true,
		"which":      true,
		"while":      true,
		"who":        true,
		"who's":      true,
		"whom":       true,
		"whose":      true,
		"why":        true,
		"why's":      true,
		"will":       true,
		"with":       true,
		"won't":      true,
		"would":      true,
		"wouldn't":   true,
		"you":        true,
		"you'd":      true,
		"you'll":     true,
		"you're":     true,
		"you've":     true,
		"your":       true,
		"yours":      true,
		"yourself":   true,
		"yourselves": true,
	}
	ng.Posts.ById = make(map[string]*Post)
	ng.Posts.ByLineNo = make(map[string]*Post)
	ng.Posts.ByShaId = make(map[string]*Post)
	ng.Posts.ByPeriod = make(map[string]*Bucket)
	ng.Posts.Spam = map[string]bool{
		"022d06d5-d618-449e-81fd-12355f80b74b@e1g2000pra.googlegroups.com":  true,
		"06946b43-14ff-4e92-8a55-2e132689fb46@a29g2000pra.googlegroups.com": true,
		"1a541972-f468-4ad9-a0a7-9eb29d111dc9@18g2000prd.googlegroups.com":  true,
		"286cba87-2bc1-4778-a302-2f309a844397@w39g2000prb.googlegroups.com": true,
		"30f94999-ca8f-4ac7-ab36-6a30d081da7e@f39g2000prb.googlegroups.com": true,
		"32b2f777-c28b-41f4-9d62-e324823e16e3@y36g2000pra.googlegroups.com": true,
		"3667db9a-a082-48a1-9507-bbe349419f86@f40g2000pri.googlegroups.com": true,
		"46ceaffa-24f6-4acc-b611-5631575ed416@t39g2000prh.googlegroups.com": true,
		"4e4585b7-e9b1-44d9-92bd-e458329161d5@r36g2000prf.googlegroups.com": true,
		"4fa888e2-9cd3-4c80-bc96-b257e17f14c0@35g2000pry.googlegroups.com":  true,
		"5f23d63e-6b76-4b95-92a2-9094239a2e38@w39g2000prb.googlegroups.com": true,
		"6129be36-9cc5-4f02-8ba7-20f84932c9d9@s4g2000yql.googlegroups.com":  true,
		"631f6fee-8a70-488d-8197-54319823ecd8@j35g2000prb.googlegroups.com": true,
		"7005d001-72d6-4a57-befb-f3ff0df675e6@j23g2000yqc.googlegroups.com": true,
		"7bb8b200-1405-4acd-b744-2bcc74b28987@v22g2000pro.googlegroups.com": true,
		"80f91ec7-ba04-46c2-9299-1819c80645f7@b31g2000prb.googlegroups.com": true,
		"813cfe70-df05-4fc8-bf48-3804c2b8fae8@w39g2000prb.googlegroups.com": true,
		"8ae4803f-1b3d-4d10-88e9-c5d0fafd0d6a@30g2000yql.googlegroups.com":  true,
		"8f53d83c-83a8-449b-8301-3dfc6c5b2cdc@v16g2000prc.googlegroups.com": true,
		"9599dfb3-e320-40aa-b33e-d8f6f1d6953b@l33g2000pri.googlegroups.com": true,
		"9e364519-e085-4006-839c-b086d3a17d53@h23g2000prf.googlegroups.com": true,
		"aaf253c0-4f8e-439f-b657-c212238c33c3@k1g2000prl.googlegroups.com":  true,
		"ad344c3a-4ee9-4f5d-88f8-5cce0e3be5ef@y32g2000prc.googlegroups.com": true,
		"b12e7025-993c-49cb-bb6d-e768d38eb527@p35g2000prm.googlegroups.com": true,
		"b34baad2-5d64-419b-8fc7-dabd60148709@n1g2000prb.googlegroups.com":  true,
		"b8e8fbed-3f74-40d0-bda8-ed28f77ec78b@a29g2000pra.googlegroups.com": true,
		"bc67562e-41e2-4d2a-937c-60c5275ea2dc@l38g2000pro.googlegroups.com": true,
		"be7f2550-8fd0-4b73-8f9c-3596cd6b6931@z7g2000prh.googlegroups.com":  true,
		"c17894f6-5720-4c9c-bb0b-569be5bcbf8e@h23g2000prf.googlegroups.com": true,
		"c2e12670-de2e-4767-82d7-fa260f20b387@k36g2000pri.googlegroups.com": true,
		"cba0e7c9-ecee-4b36-8524-61d5254ed023@q26g2000prq.googlegroups.com": true,
		"cf7c3a67-bbde-4696-829b-4fb1d2f82a33@q30g2000prq.googlegroups.com": true,
		"e094199e-4b6c-4323-8cca-94c3ee7f2f56@v11g2000prb.googlegroups.com": true,
		"e56550bf-b7a7-4c2b-86ca-c98f805f7a58@h23g2000prf.googlegroups.com": true,
		"e7502214-8acd-4a79-8386-fc5be3055805@s9g2000prg.googlegroups.com":  true,
		"e80dd951-ba26-4689-b332-d9dc09077078@d26g2000prn.googlegroups.com": true,
		"f2cab772-4740-4e4c-a894-24afbbee0b70@q8g2000prm.googlegroups.com":  true,
		"f40a0d36-b905-4e49-8a48-19e666a99e8d@a29g2000pra.googlegroups.com": true,
		"f5679ea4-556a-47a9-a76a-5da3bbcb442b@u18g2000pro.googlegroups.com": true,
		"f9f76a76-b585-42ea-b4ae-939607ad4b00@re8g2000pbc.googlegroups.com": true,
	}
	ng.Posts.Struck = make(map[string]bool)
	ng.Posts.Years = make(map[string]int)

	return ng
}

// FlagSpam will display header for suspected spam.
func (ng *NewsGroup) FlagSpam() {
	senders := map[string]bool{
		"HGHFGDS <fhfgfgg@gmail.com>": true,
		"iwcwatches5@gmail.com":       true,
	}

	for _, p := range ng.Posts.ById {
		if p.Spam {
			continue
		}
		if senders[p.Sender] {
			log.Printf("[spam] post %q\n", p.Id)
		}
	}
}

func (ng *NewsGroup) FlagStruck() {
	ids := map[string]bool{
		"sample-id": true,
	}

	for _, p := range ng.Posts.ById {
		if p.Struck {
			continue
		}
		if ids[p.Sender] {
			log.Printf("[spam] post %q\n", p.Id)
		}
	}
}

// LinkPosts links referenced and referencing posts.
func (ng *NewsGroup) LinkPosts() {
	unknownSender := "** unknown sender **"
	for _, p := range ng.Posts.ById {
		for id := range p.References {
			// when we parsed, we added a reference to the id without creating a post.
			// now we must see if that post id is in our archive. if it isn't, we
			// need to create it as a "missing" post and add it to our archive.
			xref := ng.Posts.ById[id]
			realPost := xref != nil && xref.Sender != unknownSender
			if xref == nil {
				// create it
				xref = &Post{
					Id:           id,
					ShaId:        sha1sum(id),
					Body:         "This is not the original post.\nWe were unable to locate the original in the archive.\n",
					Keys:         make(map[string][]string),
					Lines:        5,
					LineNo:       p.LineNo,
					Missing:      true,
					References:   make(map[string]*Post),
					ReferencedBy: make(map[string]*Post),
					Sender:       unknownSender,
					Subject:      "** missing post **",
				}
				// add it to the archive
				ng.Posts.ById[xref.Id] = xref
			}
			// update the link in our map
			p.References[id] = xref
			// create the back link only if the referenced post was a real post.
			if realPost {
				xref.ReferencedBy[p.Id] = p
			}
		}
	}
}

func sha1sum(s string) string {
	sum := sha1.Sum([]byte(s))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
