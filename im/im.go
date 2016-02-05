package im

/**usage
	ac := Acount{
		Uri:      "*********",
		Password: "*********",
	}

	var sender IM99U
	sender.SetAcount(&ac)

	var tos []string
	tos = append(tos, "*****")
	err := sender.SendMsg(tos, "a simple test of 99u")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
**/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var TokenValidInterval int64 = 1700

var (
	IMBASEURL     = "http://im-agent.web.sdp.101.com"
	IMTOKENPATH   = "/v0.2/api/agents/users/tokens"
	IMMESSAGEPATH = "/v0.2/api/agents/messages"
)

type Acount struct {
	Uri      string `json:"uri"`
	Password string `json:"password"`
}

type UCToken struct {
	MacAlgorithm string `json:"mac_algorithm"`
	Nonce        string `json:"nonce"`
	Mac          string `json:"mac"`
	AccessToken  string `json:"access_token"`
}

type IMapi struct {
	Filter []IMArgs `json:"filter"`
	Body   IMBody   `json:"body"`
}

type URIList struct {
	UriList []string `json:"uri_list"`
}

type GidList struct {
	GidList []string `json:"gid"`
}

type IMArgs struct {
	Name string      `json:"name"`
	Args interface{} `json:"args,omitempty"`
}

type IMBody struct {
	Content string `json:"content"`
	Flag    int    `json:"flag"`
}

type MsgResponse struct {
	Msg_id  string `json:"msg_id,omitempty"`
	Task_id string `json:"task_id,omitempty"`
	Message string `json:"message,omitempty"`
}

type IM99U struct {
	sync.Mutex
	ac         *Acount
	token      *UCToken
	lastupdate int64
}

func (this *IM99U) SetAcount(acount *Acount) {
	this.ac = acount
}

func (this *IM99U) getToken() error {
	if this.ac == nil {
		fmt.Errorf("uc acount not set")
	}
	data, err := json.Marshal(this.ac)
	if err != nil {
		return err
	}
	this.Lock()
	defer this.Unlock()
	req, _ := http.NewRequest("POST", IMBASEURL+IMTOKENPATH, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	client := http.DefaultClient
	resp, err := client.Do(req)

	token := new(UCToken)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &token)
	//	fmt.Println(string(body))
	if err != nil {
		return err
	} else {
		this.token = token
		this.lastupdate = time.Now().Unix()
		return nil
	}
}

func getIMData(tos []string, msg string) []byte {
	//根据第一个目标编号的长度判断，如果大于6位认为是群，则调用群API。小于等于6位认为是用户UID，调用个人的API
	var body IMBody
	var args IMArgs
	var api IMapi

	if len(tos[0]) <= 6 {
		var list URIList
		args.Name = "uri"
		list.UriList = tos
		args.Args = list
	}
	if len(tos[0]) > 6 {
		var list GidList
		args.Name = "gid"
		list.GidList = tos
		args.Args = list
	}
	body.Content = "Content-Type:text/plain\r\n\r\n" + msg
	body.Flag = 0
	api.Body = body
	api.Filter = append(api.Filter, args)
	data, _ := json.Marshal(&api)
	return data
}

func (this *IM99U) SendMsg(tos []string, msg string) error {
	var Users []string
	var Groups []string
	var err [2]error

	for _, id := range tos {
		if len(id) <= 6 {
			//去除空的成员
			if len(id) == 0 {
				continue
			}
			Users = append(Users, id)
		} else {
			Groups = append(Groups, id)
		}
	}
	if len(Users) != 0 {
		err[0] = this.send(Users, msg)
	}
	if len(Groups) != 0 {
		err[1] = this.send(Groups, msg)
	}
	if err[0] == nil && err[1] == nil {
		return nil
	} else {
		return fmt.Errorf("send msg error ,group user error:%s group error:%s", err[0], err[1])
	}
}

func (this *IM99U) send(tos []string, msg string) error {
	if len(tos) == 0 || (len(tos) == 1 && tos[0] == "") {
		return fmt.Errorf("tos is empty")
	}
	if msg == "" {
		return fmt.Errorf("msg is empty")
	}
	if time.Now().Unix()-this.lastupdate > TokenValidInterval {
		err := this.getToken()
		if err != nil {
			return fmt.Errorf("Get UC Token fail.with error %s", err.Error())
		}
	}
	this.Lock()
	defer this.Unlock()

	auth := fmt.Sprintf("MAC id=\"%s\",nonce=\"%s\",mac=\"%s\"", this.token.AccessToken, this.token.Nonce, this.token.Mac)
	data := bytes.NewBuffer(getIMData(tos, msg))
	req, _ := http.NewRequest("POST", IMBASEURL+IMMESSAGEPATH, data)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Authorization", auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	var msgResp MsgResponse
	body, err := ioutil.ReadAll(resp.Body)
	//	fmt.Println(string(body))
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &msgResp)
	if msgResp.Msg_id != "" || msgResp.Task_id != "" {
		return nil
	} else {
		return fmt.Errorf("Get message:%s", msgResp.Message)
	}

}
