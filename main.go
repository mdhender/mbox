package main

import (
	"flag"
	"github.com/mdhender/mbox/internal/app"
	"github.com/mdhender/mbox/internal/chunk"
	"github.com/mdhender/mbox/internal/stores/newsgroup"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	doCorpus, doSpam, showHeaders, flagSpam, flagStruck := false, false, false, false, false
	flag.BoolVar(&doCorpus, "corpus", doCorpus, "create corpus")
	flag.BoolVar(&doSpam, "spam", doCorpus, "allow spam reports")
	flag.BoolVar(&flagSpam, "flag-spam", flagSpam, "show suspected spam headers")
	flag.BoolVar(&flagStruck, "flag-struck", flagStruck, "show suspected struct headers")
	flag.BoolVar(&showHeaders, "show-headers", showHeaders, "show headers")
	flag.Parse()

	started := time.Now()
	defer func(started time.Time) {
		log.Printf("[mbox] completed in %v\n", time.Now().Sub(started))
	}(started)

	// chunks splits and cleans up the input
	chunks, err := chunk.Chunks("rec.games.pbm.mbox")
	if err != nil {
		log.Fatal(err)
	}

	ng := newsgroup.New()
	for _, ch := range chunks {
		post, err := ng.Parse(ch, doCorpus)
		if err != nil {
			log.Fatal(err)
		}
		if post.Words != nil {
			ng.Corpus.Documents[post.Id] = post.Words
		}
	}
	log.Printf("[mbox] completed parse in %v\n", time.Now().Sub(started))

	// link posts (both forwards and backwards)
	ng.LinkPosts()
	log.Printf("[mbox] completed links in %v\n", time.Now().Sub(started))

	// optional: show lists of suspect posts and quit
	if flagSpam || flagStruck {
		if flagSpam {
			ng.FlagSpam()
		}
		if flagStruck {
			ng.FlagStruck()
		}
		log.Printf("[mbox] completed flags in %v\n", time.Now().Sub(started))
		os.Exit(2)
	}

	//for _, post := range ng.SearchPosts("compliment blessed") { //  firestorm Morghoul perceval dean
	//	log.Printf("[search] post http://localhost:8080/posts/%s\n", post.ShaId)
	//}

	// we're done with the chunks
	chunks = nil

	//if post, ok := ng.Posts.ById["336E78B7.2175@earthlink.net"]; ok {
	//	log.Printf("post %q\n%q\n", post.Id, post.Body)
	//} else if post, ok := ng.Posts.ById["2s2eum$gfe@nyx10.cs.du.edu"]; ok {
	//	log.Printf("post %q\n%q\n", post.Id, post.Body)
	//}

	a, err := app.New(ng, doSpam)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[app] serving on %s\n", net.JoinHostPort(a.Host, a.Port))
	log.Fatalln(http.ListenAndServe(net.JoinHostPort(a.Host, a.Port), a.Router))
}
