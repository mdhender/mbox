package newsgroup

import (
	"crypto/sha1"
	"encoding/base64"
	"github.com/mdhender/mbox/internal/chunk"
	"log"
)

type NewsGroup struct {
	Corpus struct {
		// Documents is a list of all posts with word counts
		Documents map[string]map[string]int
		// Index is a map of word to posts that contain the word
		Index map[string][]*Post
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
	ng.Corpus.Documents = make(map[string]map[string]int)
	ng.Corpus.Index = make(map[string][]*Post)
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
		"01bdec34$586dfe60$b0e00b80@hp-pavilion":                              true,
		"01bdeca0$b784d420$27edb5cf@conquest":                                 true,
		"022d06d5-d618-449e-81fd-12355f80b74b@e1g2000pra.googlegroups.com":    true,
		"06946b43-14ff-4e92-8a55-2e132689fb46@a29g2000pra.googlegroups.com":   true,
		"19961231183000.NAA16482@ladder01.news.aol.com":                       true,
		"1997Jan3.105704.18682@leeds.ac.uk":                                   true,
		"1998042812192800.IAA12113@ladder03.news.aol.com":                     true,
		"1998042923214600.TAA26610@ladder03.news.aol.com":                     true,
		"1998050416352200.MAA18095@ladder01.news.aol.com":                     true,
		"1998050716290500.MAA17203@ladder01.news.aol.com":                     true,
		"1998051121370700.RAA11554@ladder03.news.aol.com":                     true,
		"1998051122523200.SAA18140@ladder03.news.aol.com":                     true,
		"1998051122555500.SAA18600@ladder03.news.aol.com":                     true,
		"1998051218042500.OAA18281@ladder03.news.aol.com":                     true,
		"1998051220481900.QAA08372@ladder01.news.aol.com":                     true,
		"1998051220484900.QAA08445@ladder01.news.aol.com":                     true,
		"1998051220493700.QAA04317@ladder03.news.aol.com":                     true,
		"1998051220495000.QAA08552@ladder01.news.aol.com":                     true,
		"1998051220500300.QAA04372@ladder03.news.aol.com":                     true,
		"1998051220501200.QAA08606@ladder01.news.aol.com":                     true,
		"1998051220502500.QAA04427@ladder03.news.aol.com":                     true,
		"1998051220503900.QAA04460@ladder03.news.aol.com":                     true,
		"1998051220510000.QAA08713@ladder01.news.aol.com":                     true,
		"1998051307314200.DAA17182@ladder01.news.aol.com":                     true,
		"1998051307321900.DAA12969@ladder03.news.aol.com":                     true,
		"1998051307332500.DAA17236@ladder01.news.aol.com":                     true,
		"1998051315340100.LAA04793@ladder03.news.aol.com":                     true,
		"1998051420134400.QAA27279@ladder01.news.aol.com":                     true,
		"1998051423084900.TAA16506@ladder01.news.aol.com":                     true,
		"1998051423093000.TAA12333@ladder03.news.aol.com":                     true,
		"1998051518544900.OAA13161@ladder01.news.aol.com":                     true,
		"1998051723504500.TAA17299@ladder03.news.aol.com":                     true,
		"1998090716342000.MAA04316@ladder03.news.aol.com":                     true,
		"1998090919153900.PAA21789@ladder03.news.aol.com":                     true,
		"1998090919154600.PAA18712@ladder01.news.aol.com":                     true,
		"1998091011415401.HAA06391@ladder03.news.aol.com":                     true,
		"1998091011482100.HAA06751@ladder03.news.aol.com":                     true,
		"1998091011483600.HAA03917@ladder01.news.aol.com":                     true,
		"1998091115464600.LAA17094@ladder03.news.aol.com":                     true,
		"1998091115465700.LAA14723@ladder01.news.aol.com":                     true,
		"1998091115470600.LAA17131@ladder03.news.aol.com":                     true,
		"19980917181751.17754.00000713@ng08.aol.com":                          true,
		"19980917181801.17754.00000714@ng08.aol.com":                          true,
		"19980917181813.17754.00000715@ng08.aol.com":                          true,
		"19980918174826.22787.00002467@ng67.aol.com":                          true,
		"19980920163307.06732.00002371@ng-fa2.aol.com":                        true,
		"19980925140621.11768.00001557@ng-fa1.aol.com":                        true,
		"19980925195003.08090.00001778@ng110.aol.com":                         true,
		"19980926124214.04111.00002146@ng25.aol.com":                          true,
		"19980927143426.09342.00002857@ng109.aol.com":                         true,
		"19980928142610.11772.00003088@ng-fa1.aol.com":                        true,
		"19980929120604.15632.00003749@ng43.aol.com":                          true,
		"19980929170015.18951.00003994@ng86.aol.com":                          true,
		"19980930091422.05669.00004073@ng112.aol.com":                         true,
		"19980930091744.05669.00004076@ng112.aol.com":                         true,
		"19981001160051.28440.00005170@ng37.aol.com":                          true,
		"19981001160327.28440.00005182@ng37.aol.com":                          true,
		"19981001160344.28440.00005183@ng37.aol.com":                          true,
		"19981009161035.10310.00008566@ng16.aol.com":                          true,
		"19990113174600.02507.00000303@ng101.aol.com":                         true,
		"1a541972-f468-4ad9-a0a7-9eb29d111dc9@18g2000prd.googlegroups.com":    true,
		"20000330192429.24408.00000084@nso-cm.aol.com":                        true,
		"286cba87-2bc1-4778-a302-2f309a844397@w39g2000prb.googlegroups.com":   true,
		"30f94999-ca8f-4ac7-ab36-6a30d081da7e@f39g2000prb.googlegroups.com":   true,
		"32C46F0B.5A1E@public.srce.hr":                                        true,
		"32C52DD5.4658@ic.mankato.mn.us":                                      true,
		"32C6F863.1FC8@io.com":                                                true,
		"32C6F934.1A32@io.com":                                                true,
		"32CDEF8B.348C@ic.mankato.mn.us":                                      true,
		"32b2f777-c28b-41f4-9d62-e324823e16e3@y36g2000pra.googlegroups.com":   true,
		"33369fcd.9078223@news.dave-world.net":                                true,
		"3547C669.446@direct.ca":                                              true,
		"3548A6E1.1FB9DA2F@ibm.net":                                           true,
		"3548B93D.8E69F1F0@ibm.net":                                           true,
		"3551ecd5.0@news.saqnet.co.uk":                                        true,
		"35C53BE2.63564842@igergy.cl":                                         true,
		"35F2013C.EFD4763D@the-isles.demon.co.uk":                             true,
		"35F70D59.52CB6BC@usa.net":                                            true,
		"35F7F755.41C6@mitre.org":                                             true,
		"35F8F6D6.7A4C45D@iaehv.nl":                                           true,
		"35FD330B.5DDE@Toronto.BCSC.ON.Bell.CA":                               true,
		"3609A3BF.3F83A14F@usa.net":                                           true,
		"3611E873.C2D5BC58@ferndown.ate.slb.com":                              true,
		"3611F0E7.456A@signal.dera.gov.uk":                                    true,
		"36120101.45B4@signal.dera.gov.uk":                                    true,
		"3667db9a-a082-48a1-9507-bbe349419f86@f40g2000pri.googlegroups.com":   true,
		"38e6494a.826543@news.cwcom.net":                                      true,
		"38e7b432.614299@news.cwcom.net":                                      true,
		"3B9C0486.DA57448D@fuse.net":                                          true,
		"3B9C2063.2B50F7D1@fuse.net":                                          true,
		"3B9C2B99.BEDBB042@sympatico.ca":                                      true,
		"3EC84E43.5C0BDF47@worldnet.att.net":                                  true,
		"3c056ec1_8@Usenet.com":                                               true,
		"3c05f3a8.56228362@Internet":                                          true,
		"46ceaffa-24f6-4acc-b611-5631575ed416@t39g2000prh.googlegroups.com":   true,
		"4B0BF17FDCD3039B.9BD935610AD37D23.075A04FA67373488@lp.airnews.net":   true,
		"4HAm7.6704$d86.549992@newsread1.prod.itd.earthlink.net":              true,
		"4Nxm7.121$d4.11916@sc0101.promedia.net":                              true,
		"4QKXlGAJlfU1Ewad@homeway.demon.co.uk":                                true,
		"4a21e522.0406180132.7217485c@posting.google.com":                     true,
		"4e4585b7-e9b1-44d9-92bd-e458329161d5@r36g2000prf.googlegroups.com":   true,
		"4fa888e2-9cd3-4c80-bc96-b257e17f14c0@35g2000pry.googlegroups.com":    true,
		"4hcv1s$e3a@thrush.sover.net":                                         true,
		"4hcv1v$e3a@thrush.sover.net":                                         true,
		"5483b912.0111282110.2985143e@posting.google.com":                     true,
		"5483b912.0111291720.799fc3e2@posting.google.com":                     true,
		"5913fd86.1fcb2cb1@host-69-48-73-244.roc.choiceone.net":               true,
		"5a29pt$jad2@mars.online.uleth.ca":                                    true,
		"5a3us6$h59@leaphome.demon.co.uk":                                     true,
		"5a4t3v$5b2@nr1.vancouver.istar.net":                                  true,
		"5abr2d$lgb@inet-nntp-gw-1.us.oracle.com":                             true,
		"5f23d63e-6b76-4b95-92a2-9094239a2e38@w39g2000prb.googlegroups.com":   true,
		"5ont4t$st1$2320@news.internetmci.com":                                true,
		"5p6mgg$qa3$224@news.internetmci.com":                                 true,
		"6129be36-9cc5-4f02-8ba7-20f84932c9d9@s4g2000yql.googlegroups.com":    true,
		"62c72772.0406210202.140e6dab@posting.google.com":                     true,
		"631f6fee-8a70-488d-8197-54319823ecd8@j35g2000prb.googlegroups.com":   true,
		"6i9gi6$3v4@ask.diku.dk":                                              true,
		"6i9v5a$6aq$2@nclient3-gui.server.virgin.net":                         true,
		"6t1a4e$pdv$1@zeus.tcp.net.uk":                                        true,
		"6t6mfg$2tu@news1.newsguy.com":                                        true,
		"6tc513$64m$1@zeus.tcp.net.uk":                                        true,
		"6urk21$8e1@news1.newsguy.com":                                        true,
		"6uuf2t$n2c$1@nnrp1.dejanews.com":                                     true,
		"6vhb06$rcu$1@usenet40.supernews.com":                                 true,
		"7005d001-72d6-4a57-befb-f3ff0df675e6@j23g2000yqc.googlegroups.com":   true,
		"7bb8b200-1405-4acd-b744-2bcc74b28987@v22g2000pro.googlegroups.com":   true,
		"7yt6ni$4eq$1@news.arcor-ip.net":                                      true,
		"80f91ec7-ba04-46c2-9299-1819c80645f7@b31g2000prb.googlegroups.com":   true,
		"813cfe70-df05-4fc8-bf48-3804c2b8fae8@w39g2000prb.googlegroups.com":   true,
		"894557509.12147.1.nnrp-02.9e989556@news.demon.co.uk":                 true,
		"894965999.23367.0.nnrp-01.9e989556@news.demon.co.uk":                 true,
		"894994329.10648.0.nnrp-04.d4e42dee@news.demon.co.uk":                 true,
		"8ae4803f-1b3d-4d10-88e9-c5d0fafd0d6a@30g2000yql.googlegroups.com":    true,
		"8c4dv9$36e$1@plutonium.btinternet.com":                               true,
		"8c76lu$lq7$1@uranium.btinternet.com":                                 true,
		"8cdav1$pms$1@neptunium.btinternet.com":                               true,
		"8f53d83c-83a8-449b-8301-3dfc6c5b2cdc@v16g2000prc.googlegroups.com":   true,
		"905422914.11852.6.nnrp-09.9e989556@news.demon.co.uk":                 true,
		"907093147.17060.0.nnrp-01.9e989556@news.demon.co.uk":                 true,
		"907114511.7052.0.nnrp-09.d4e4933e@news.demon.co.uk":                  true,
		"907148984.18391.2.nnrp-04.9e989556@news.demon.co.uk":                 true,
		"907178501.3055.1.nnrp-07.9e989556@news.demon.co.uk":                  true,
		"907194001.1187.0.nnrp-05.d4e4933e@news.demon.co.uk":                  true,
		"940b809b-0769-422d-a0b3-b8722df18a57@s12g2000prg.googlegroups.com":   true,
		"954880066.26619.0.nnrp-14.9e989556@news.demon.co.uk":                 true,
		"9599dfb3-e320-40aa-b33e-d8f6f1d6953b@l33g2000pri.googlegroups.com":   true,
		"9E4BCCD1516658E38274ADCBC8@news.arcor-ip.net":                        true,
		"9e364519-e085-4006-839c-b086d3a17d53@h23g2000prf.googlegroups.com":   true,
		"9fWmmNAqigA2EwmU@homeway.demon.co.uk":                                true,
		"A0FZcJAueLS1EwPW@homeway.demon.co.uk":                                true,
		"CcjF4.1852$Nc2.38494@news3.cableinet.net":                            true,
		"E7AoPAAGYE91EwH3@homeway.demon.co.uk":                                true,
		"F4QYPAAlZ3V1EwXd@homeway.demon.co.uk":                                true,
		"F9fI+DATZSn2EwbA@homeway.demon.co.uk":                                true,
		"Ge5F4.803$Nc2.12753@news3.cableinet.net":                             true,
		"Ht1nrJAA1fS1Ewvw@mhairi.demon.co.uk":                                 true,
		"MANaRDA8QfO1Ewt1@homeway.demon.co.uk":                                true,
		"MPG.1081b6cca3e516b49896a8@news.cport.com":                           true,
		"MPG.108217923e1b54b09896ac@news.cport.com":                           true,
		"OABQEzbe9GA.319@ntawwabp.compuserve.com":                             true,
		"Pine.A32.3.93.961228133812.31714A-100000@srv1.freenet.calgary.ab.ca": true,
		"Pine.GSO.3.96.980430092913.28423B-100000@dante":                      true,
		"RCIF4.3671$Nc2.93873@news3.cableinet.net":                            true,
		"RP$OMHA5NtU1Ew$J@tincat.demon.co.uk":                                 true,
		"UZvjd.8167$O11.2158@newsread3.news.pas.earthlink.net":                true,
		"aaf253c0-4f8e-439f-b657-c212238c33c3@k1g2000prl.googlegroups.com":    true,
		"ab0800011439471637@4ax.com":                                          true,
		"ab512dd8.0406170239.39303466@posting.google.com":                     true,
		"acb07122.0409230135.7826ed8e@posting.google.com":                     true,
		"ad0800011215340695@4ax.com":                                          true,
		"ad344c3a-4ee9-4f5d-88f8-5cce0e3be5ef@y32g2000prc.googlegroups.com":   true,
		"ama029$bhs$21980@news1.kornet.net":                                   true,
		"b12e7025-993c-49cb-bb6d-e768d38eb527@p35g2000prm.googlegroups.com":   true,
		"b34baad2-5d64-419b-8fc7-dabd60148709@n1g2000prb.googlegroups.com":    true,
		"b8e8fbed-3f74-40d0-bda8-ed28f77ec78b@a29g2000pra.googlegroups.com":   true,
		"bc0711311144590514@4ax.com":                                          true,
		"bc67562e-41e2-4d2a-937c-60c5275ea2dc@l38g2000pro.googlegroups.com":   true,
		"bckkl7$jiiq2$1@ID-158456.news.dfncis.de":                             true,
		"be7f2550-8fd0-4b73-8f9c-3596cd6b6931@z7g2000prh.googlegroups.com":    true,
		"bjm10-0210981232110001@potato.cit.cornell.edu":                       true,
		"bk5F4.811$Nc2.12884@news3.cableinet.net":                             true,
		"c17894f6-5720-4c9c-bb0b-569be5bcbf8e@h23g2000prf.googlegroups.com":   true,
		"c2e12670-de2e-4767-82d7-fa260f20b387@k36g2000pri.googlegroups.com":   true,
		"cba0e7c9-ecee-4b36-8524-61d5254ed023@q26g2000prq.googlegroups.com":   true,
		"cf7c3a67-bbde-4696-829b-4fb1d2f82a33@q30g2000prq.googlegroups.com":   true,
		"ch0711310845523807@4ax.com":                                          true,
		"cv0711311207246699@4ax.com":                                          true,
		"d32dcf62.d0cc66ef@lex5z7z2ghz.com":                                   true,
		"dm0800010947291306@4ax.com":                                          true,
		"e094199e-4b6c-4323-8cca-94c3ee7f2f56@v11g2000prb.googlegroups.com":   true,
		"e0W2ViK79GA.49@nih2naaa.prod2.compuserve.com":                        true,
		"e56550bf-b7a7-4c2b-86ca-c98f805f7a58@h23g2000prf.googlegroups.com":   true,
		"e7502214-8acd-4a79-8386-fc5be3055805@s9g2000prg.googlegroups.com":    true,
		"e80dd951-ba26-4689-b332-d9dc09077078@d26g2000prn.googlegroups.com":   true,
		"ee84fb56-9191-4dba-a487-99d186ef66f7@x69g2000hsx.googlegroups.com":   true,
		"etcld.63$%Z.61@fe37.usenetserver.com":                                true,
		"f2cab772-4740-4e4c-a894-24afbbee0b70@q8g2000prm.googlegroups.com":    true,
		"f40a0d36-b905-4e49-8a48-19e666a99e8d@a29g2000pra.googlegroups.com":   true,
		"f5679ea4-556a-47a9-a76a-5da3bbcb442b@u18g2000pro.googlegroups.com":   true,
		"f9f76a76-b585-42ea-b4ae-939607ad4b00@re8g2000pbc.googlegroups.com":   true,
		"fk0711311017246578@4ax.com":                                          true,
		"gh0800020811258065@4ax.com":                                          true,
		"hk0711260923269136@4ax.com":                                          true,
		"ix0711310905256481@4ax.com":                                          true,
		"jw0711310951351952@4ax.com":                                          true,
		"m07012716453074@4ax.com":                                             true,
		"m07072011255906@4ax.com":                                             true,
		"m07072013252228@4ax.com":                                             true,
		"m07072014140368@4ax.com":                                             true,
		"m07072014410343@4ax.com":                                             true,
		"m07072015063897@4ax.com":                                             true,
		"m07072015274460@4ax.com":                                             true,
		"m07072016081847@4ax.com":                                             true,
		"m07072017225489@4ax.com":                                             true,
		"m07072018450193@4ax.com":                                             true,
		"m07072020005868@4ax.com":                                             true,
		"m07072506163270@4ax.com":                                             true,
		"m07072507505830@4ax.com":                                             true,
		"m07072508165406@4ax.com":                                             true,
		"m07072510231745@4ax.com":                                             true,
		"m07072512080703@4ax.com":                                             true,
		"m07072512481436@4ax.com":                                             true,
		"m07101111325958@4ax.com":                                             true,
		"m07101113314996@4ax.com":                                             true,
		"m07101114374997@4ax.com":                                             true,
		"m07101709073990@4ax.com":                                             true,
		"m07101716031378@4ax.com":                                             true,
		"m07101717415628@4ax.com":                                             true,
		"m07101722062936@4ax.com":                                             true,
		"m07101803261011@4ax.com":                                             true,
		"m07101811450285@4ax.com":                                             true,
		"m07101813355211@4ax.com":                                             true,
		"m07101816294145@4ax.com":                                             true,
		"m07101822244674@4ax.com":                                             true,
		"m07102407422631@4ax.com":                                             true,
		"m07102408113148@4ax.com":                                             true,
		"m07102408384904@4ax.com":                                             true,
		"m07102409110540@4ax.com":                                             true,
		"m07102410134981@4ax.com":                                             true,
		"m07110113032381@4ax.com":                                             true,
		"m07111610471217@4ax.com":                                             true,
		"m07111611470128@4ax.com":                                             true,
		"m07111612141308@4ax.com":                                             true,
		"m07111612413850@4ax.com":                                             true,
		"m07111613055896@4ax.com":                                             true,
		"m07111613291228@4ax.com":                                             true,
		"m07111613560999@4ax.com":                                             true,
		"m07111614232578@4ax.com":                                             true,
		"m07111614522962@4ax.com":                                             true,
		"m07111615201713@4ax.com":                                             true,
		"m07111615502245@4ax.com":                                             true,
		"m07111616195337@4ax.com":                                             true,
		"m07111616511608@4ax.com":                                             true,
		"m07111617185399@4ax.com":                                             true,
		"m07111617563406@4ax.com":                                             true,
		"mp0711310825190405@4ax.com":                                          true,
		"ni0800011259076974@4ax.com":                                          true,
		"py0800011607374526@4ax.com":                                          true,
		"rH9jQPA80Y+1EwYI@homeway.demon.co.uk":                                true,
		"sd0711310927478889@4ax.com":                                          true,
		"sfGdnTGFeZSRNBPcRVn-tw@comcast.com":                                  true,
		"srt.845656157@sun-dimas":                                             true,
		"tSu_.52$E6.317811@ptah.visi.com":                                     true,
		"tUgHa.4600$LP.211@newsfep4-winn.server.ntli.net":                     true,
		"uT20xx$69GA.312@ntdwwaaw.compuserve.com":                             true,
		"vlwN7.241259$W8.8439676@bgtnsc04-news.ops.worldnet.att.net":          true,
		"ykkuazmmqg37cgkv51htm3w0bef1h@4ax.com":                               true,
		"za0711311058133702@4ax.com":                                          true,
		"zwxVdAAnvfE2EwUs@homeway.demon.co.uk":                                true,
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

func (ng *NewsGroup) SearchPosts(input string) map[string]*Post {
	posts := make(map[string]*Post)
	for n, word := range chunk.Tokenize([]byte(input), ng.Corpus.StopWords) {
		// set is the set of documents containing this word
		set := make(map[string]*Post)
		if documents, ok := ng.Corpus.Index[string(word)]; ok {
			for _, post := range documents {
				set[post.ShaId] = post
			}
		}
		if n == 0 {
			posts = set
		} else {
			// we want the intersection of current posts and this set
			intersection := make(map[string]*Post)
			for _, post := range set {
				if _, ok := posts[post.ShaId]; ok {
					intersection[post.ShaId] = post
				}
			}
			// result is the intersection
			posts = intersection
		}
	}
	return posts
}
