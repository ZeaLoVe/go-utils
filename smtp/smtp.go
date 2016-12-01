package smtp

import (
	"crypto/tls"
	"fmt"
	"net"
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

func CopyValideInfo(list []string) []string {
	var ret []string
	for _, val := range list {
		if val != "" {
			ret = append(ret, val)
		}
	}
	return ret
}

/*修改原来的sendmail函数*/
func MySendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		host, _, _ := net.SplitHostPort(addr)
		config := &tls.Config{InsecureSkipVerify: true, ServerName: host}
		if err = c.StartTLS(config); err != nil {
			return err
		}
	}
	if a != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(a); err != nil {
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
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
	err := MySendMail(this.ac.Server+SMTPPORTSTRING, this.auth, this.ac.User, CopyValideInfo(strings.Split(mail.Tos, ",")), []byte(mail.String()))
	if err != nil {
		return err
	}
	return nil
}
