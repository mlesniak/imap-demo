package main

import (
	"fmt"
	"os"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	color "github.com/logrusorgru/aurora"
)

type config struct {
	Username      string
	Password      string
	NumberMessage uint32
}

var credentials = config{os.Getenv("IMAP_USERNAME"), os.Getenv("IMAP_PASSWORD"), 15}

func main() {
	// Connect to server
	c, _ := client.DialTLS("mail.micromata.de:993", nil)
	defer c.Logout()
	err := c.Login(credentials.Username, credentials.Password)
	if err != nil {
		fmt.Println(color.Red("Unable to login, is IMAP_USERNAME/PASSWORD set?"))
		return
	}

	// Select INBOX
	mbox, _ := c.Select("INBOX", false)

	to := mbox.Messages
	from := mbox.Messages - credentials.NumberMessage
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	go func() {
		c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	for msg := range messages {
		envelope := msg.Envelope
		from := (*envelope.Sender[0]).PersonalName
		fromLen := 20
		if len(from) < 20 {
			fromLen = len(from)
		}
		subject := envelope.Subject
		subLen := 60
		if len(subject) < 60 {
			subLen = len(subject)
		}
		fmt.Printf("%s: %s\n", color.Magenta(from[:fromLen]), subject[:subLen])
	}
}
