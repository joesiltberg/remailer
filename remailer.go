package main

import (
	"flag"
	"fmt"
	"io"
	"net/mail"
	"os"
	"strings"
)

var headersToCopy = []string{
	"Mime-Version",
	"Date",
	"Content-Type",
	"Content-Transfer-Encoding",
}

const exitFailure = 1

func main() {
	var bounceAddressp = flag.String("sender", "", "The sender address of the new message (From:)")
	var remailerNamep = flag.String("name", "", "Full name of the sender")
	var recipientsp = flag.String("recipients", "", "Space separated list of recipients (To:)")
	var stripp = flag.String("strip", "", "If given, the first instance of this string is removed from the subject line")

	flag.Parse()

	bounceAddress := *bounceAddressp
	remailerName := *remailerNamep
	recipients := *recipientsp
	strip := *stripp

	msg, err := mail.ReadMessage(os.Stdin)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read mail: %v", err)
		os.Exit(exitFailure)
	}

	for _, s := range headersToCopy {
		for i := range msg.Header[s] {
			fmt.Printf("%s: %s\n", s, msg.Header[s][i])
		}
	}

	if replyTo, ok := msg.Header["Reply-To"]; ok && len(replyTo) > 0 {
		fmt.Printf("Reply-To: %s\n", replyTo[0])
	} else if from, ok := msg.Header["From"]; ok && len(from) > 0 {
		fmt.Printf("Reply-To: %s\n", from[0])
	}
	fmt.Printf("From: \"%s\" <%s>\n", remailerName, bounceAddress)

	if strip != "" {
		if subject, ok := msg.Header["Subject"]; ok && len(subject) > 0 {
			fmt.Printf("Subject: %s\n", strings.Replace(subject[0], strip, "", 1))
		}
	}

	recipientsCommaSeparated := strings.Join(strings.Fields(recipients), ", ")

	fmt.Printf("To: %s\n", recipientsCommaSeparated)

	io.Copy(os.Stdout, msg.Body)
}
