package main

import (
	"strings"
	"fmt"
	"io/ioutil"
	"github.com/VishalTanwani/GolangWebApp/internal/modals"
	"github.com/xhit/go-simple-mail/v2"
	"log"
	"time"
)

func listenForMail() {
	for {
		msg := <-app.MailChan
		sendMail(msg)
	}
}

func sendMail(m modals.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		log.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data,err := ioutil.ReadFile(fmt.Sprintf("./email-templates/%s",m.Template))
		if err!=nil{
			log.Println(err)
		}
		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate,"[%body%]",m.Content,1)
		email.SetBody(mail.TextHTML, msgToSend)
	}

	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("MailSend")
	}

}
