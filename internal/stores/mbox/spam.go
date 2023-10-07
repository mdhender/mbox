package mbox

import "log"

func (mb *MailBox) FlagSpam() {
	// code to display header for suspected spam. should be commented out.
	for _, msg := range mb.Messages {
		if msg.Spam {
			continue
		}
		if msg.Header.From == "HGHFGDS <fhfgfgg@gmail.com>" {
			log.Printf("[spam] %q\n", msg.Header.Id)
		} else if msg.Header.From == "iwcwatches5@gmail.com" {
			log.Printf("[spam] %q\n", msg.Header.Id)
		}
	}
}

var spam = map[string]bool{
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
