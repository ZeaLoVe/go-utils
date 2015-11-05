package smtp

import (
	"fmt"
	"net/smtp"
	"strings"
)

var SMTPPORTSTRING = ":25"

type Mail struct {
	Tos     string `json:"tos"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func (this *Mail) String() string {
	return fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\n\r\n%s\r\r",
		this.Tos,
		this.Subject,
		this.Content,
	)
}

type MailAcount struct {
	User     string
	Password string
	Server   string
}

type SmtpMailSender struct {
	ac *MailAcount
}

func (this *SmtpMailSender) SetMailAcount(acount *MailAcount) error {
	if acount == nil || acount.Password == "" || acount.User == "" || acount.Password == "" {
		return fmt.Errorf("set acount error")
	}
	this.ac = acount
	return nil
}

func (this *SmtpMailSender) SendMail(mail *Mail) error {
	auth := smtp.PlainAuth("", this.ac.User, this.ac.Password, this.ac.Server)
	err := smtp.SendMail(this.ac.Server+SMTPPORTSTRING, auth, this.ac.User, strings.Split(mail.Tos, ","), []byte(mail.String()))
	if err != nil {
		//		fmt.Println(err)
		return err
	}
	return nil
}
