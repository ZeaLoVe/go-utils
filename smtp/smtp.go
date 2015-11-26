package smtp

import (
	"fmt"
	"net/smtp"
	"strings"
	"sync"
)

var SMTPPORTSTRING = ":25"

type Mail struct {
	Tos     string `json:"tos"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func (this *Mail) String() string {
	return fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nContent-Type: text/plain;charset=UTF-8\r\n\r\n%s\r\r",
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
	lock *sync.Mutex
	ac   *MailAcount
	auth smtp.Auth
}

func (this *SmtpMailSender) SetMailAcount(acount *MailAcount) error {
	if acount == nil || acount.Password == "" || acount.User == "" || acount.Server == "" {
		return fmt.Errorf("set acount error")
	}
	this.lock = &sync.Mutex{}
	this.ac = acount
	this.auth = smtp.PlainAuth("", this.ac.User, this.ac.Password, this.ac.Server)
	return nil
}

func (this *SmtpMailSender) SendMail(mail *Mail) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.auth == nil {
		return fmt.Errorf("auth not corret")
	}
	if mail.Tos == "" || mail.Content == "" || mail.Subject == "" {
		return fmt.Errorf("mail infomation lack")
	}
	err := smtp.SendMail(this.ac.Server+SMTPPORTSTRING, this.auth, this.ac.User, strings.Split(mail.Tos, ","), []byte(mail.String()))
	if err != nil {
		return err
	}
	return nil
}
