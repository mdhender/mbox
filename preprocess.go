package main

import (
	"bytes"
	"regexp"
)

// PreProcess is a set of hacks to "fix" problematic inputs.
// It is written specifically for the `rec.games.pbm` mbox file.
func preprocess(lines [][]byte) [][]byte {
	// every message starts with a blank line followed by a From with ID.
	refTrailingSpace := regexp.MustCompile("^References: .* $")

	for i, line := range lines {
		if bytes.HasPrefix(line, []byte("Date: ")) {
			if bytes.Equal(line, []byte("Date: 11 Sep 93 12:58:28 -500")) {
				lines[i] = []byte("Date: 11 Sep 93 12:58:28 -0500")
			} else if bytes.Equal(line, []byte("Date: 11 Sep 93 23:10:45 -500")) {
				lines[i] = []byte("Date: 11 Sep 93 23:10:45 -0500")
			} else if bytes.Equal(line, []byte("Date: Wed, 12 Oct 1994 09:35:51 Central")) {
				lines[i] = []byte("Date: Wed, 12 Oct 1994 09:35:51 CST")
			} else if bytes.Equal(line, []byte("Date: Thu, 02 Dec 93 19:50:54 est")) {
				lines[i] = []byte("Date: Thu, 02 Dec 93 19:50:54 EST")
			} else if bytes.Equal(line, []byte("Date: Tue, 15 Jun 93 15:10:37 T-1")) {
				lines[i] = []byte("Date: Tue, 15 Jun 93 15:10:37 -0100")
			}
			continue
		}
		if bytes.HasPrefix(line, []byte("References: ")) {
			if refTrailingSpace.Find(line) != nil {
				lines[i] = bytes.TrimSpace(line)
			}
			if bytes.Equal(line, []byte("References: <")) {
				lines[i] = []byte("References: <missing-reference-id>")
			} else if bytes.Equal(line, []byte("References: C0GzED.A2u@news.cso.uiuc.edu> <1829@idacrd.UUCP> <1ii5rfINNc2q@darkstar.UCSC.EDU")) {
				lines[i] = []byte("References: <C0GzED.A2u@news.cso.uiuc.edu> <1829@idacrd.UUCP> <1ii5rfINNc2q@darkstar.UCSC.EDU>")
			} else if bytes.Equal(line, []byte("References: RSI Customer Service")) {
				lines[i] = []byte("References: <RSI-Customer-Service>")
			} else if bytes.Equal(line, []byte("References: <1991Apr13.030312.7999@vax1.tcd.ie}")) {
				lines[i] = []byte("References: <1991Apr13.030312.7999@vax1.tcd.ie>")
			} else if bytes.Equal(line, []byte("References: <1991Nov12.183857.24316@newcastle.ac.uk> <1991Nov18.011915.40")) {
				lines[i] = []byte("References: <1991Nov12.183857.24316@newcastle.ac.uk> <1991Nov18.011915.408@bradley.bradley.edu>")
			} else if bytes.Equal(line, []byte("References: <1992Mar21.004047.17322@erg.sri.com>> <18182@ector.cs.purdue.edu> <1992Mar21.213430.8671@daimi.aau.dk")) {
				//lines[i] = []byte("References: <1992Mar21.004047.17322@erg.sri.com> <18182@ector.cs.purdue.edu> <1992Mar21.213430.8671@daimi.aau.dk>")
			} else if bytes.Equal(line, []byte("References: <1993Feb1.162305.16901@magnus.acs.ohio-state.edu> <1kjon1INN81d@bre")) {
				lines[i] = []byte("References: <1993Feb1.162305.16901@magnus.acs.ohio-state.edu> <1kjon1INN81d@bredbeddle.cs.purdue.edu>")
			} else if bytes.Equal(line, []byte("References: <8fJ=SMe00WBLE7En4P@andrew.cmu.edu> <21390@ucdavis.ucdavis.edu> <8f")) {
				lines[i] = []byte("References: <8fJ=SMe00WBLE7En4P@andrew.cmu.edu> <21390@ucdavis.ucdavis.edu> <invalid-reference-id>")
			} else if bytes.Equal(line, []byte("References: <C1tyDE.EI9@inews.Intel.COM> <16B69C2D4.X049RH@tamvm1.tamu.edu> <19")) {
				lines[i] = []byte("References: <C1tyDE.EI9@inews.Intel.COM> <16B69C2D4.X049RH@tamvm1.tamu.edu> <1993Feb4.044100.17009@midway.uchicago.edu>")
			}
		}
	}

	return lines
}
