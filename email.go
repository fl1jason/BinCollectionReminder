package main

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
)

type Email struct {
	Recipient string
	Template  string
	Image     string
	BinType   string
}

func send(e Email) {

	// Sender data.
	from := "samjas73@gmail.com"
	password := "SamJas73!"

	// Receiver email address.
	to := []string{
		e.Recipient,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    "Jason",
		Message: "Hello World",
	})

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}
