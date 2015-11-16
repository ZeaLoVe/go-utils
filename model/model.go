package model

import (
	"fmt"
)

//所有的TOS都是通过 ','逗号分隔的

type Sms struct {
	Tos     string `json:"tos"`
	Content string `json:"content"`
}

//IM sms,Tos need to be 99u or any company im acount
type IMSms struct {
	Tos     string `json:"tos"`
	Content string `json:"content"`
}

//phone message,Tos means phone number
type Phone struct {
	Tos     string `json:"tos"`
	Content string `json:"content"`
}

type Mail struct {
	Tos     string `json:"tos"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func (this *Sms) String() string {
	return fmt.Sprintf(
		"<Tos:%s, Content:%s>",
		this.Tos,
		this.Content,
	)
}

func (this *IMSms) String() string {
	return fmt.Sprintf(
		"<Tos:%s, Content:%s>",
		this.Tos,
		this.Content,
	)
}

func (this *Phone) String() string {
	return fmt.Sprintf(
		"<Tos:%s, Content:%s>",
		this.Tos,
		this.Content,
	)
}

//func (this *Mail) String() string {
//	return fmt.Sprintf(
//		"<Tos:%s, Subject:%s, Content:%s>",
//		this.Tos,
//		this.Subject,
//		this.Content,
//	)
//}

//不希望看到content那么多行，简单的输出一行就好
func (this *Mail) String() string {
	return fmt.Sprintf(
		"<Tos:%s, Subject:%s>",
		this.Tos,
		this.Subject,
	)
}
