package im

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

type IMList struct {
	UriList []string `json:"uri_list"`
}

type IMArgs struct {
	Name string `json:"name"`
	Args IMList `json:"args"`
}

type IMBody struct {
	Content string `json:"content"`
	Flag    int    `json:"flag"`
}

type MsgResponse struct {
	Msg_id  string `json:"msg_id,omitempty"`
	Message string `json:"message,omitempty"`
}

type IM99U struct {
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
	var list IMList
	var body IMBody
	var args IMArgs
	var api IMapi
	list.UriList = tos
	args.Name = "uri"
	args.Args = list
	body.Content = "Content-Type:text/plain\r\n\r\n" + msg
	body.Flag = 0
	api.Body = body
	api.Filter = append(api.Filter, args)
	data, _ := json.Marshal(&api)
	return data
}

func (this *IM99U) SendMsg(tos []string, msg string) error {
	if time.Now().Unix()-this.lastupdate > TokenValidInterval {
		fmt.Println("Get token called")
		err := this.getToken()
		if err != nil {
			return fmt.Errorf("Get UC Token fail.with error %s", err.Error())
		}
	}
	auth := fmt.Sprintf("MAC id=\"%s\",nonce=\"%s\",mac=\"%s\"", this.token.AccessToken, this.token.Nonce, this.token.Mac)
	data := bytes.NewBuffer(getIMData(tos, msg))
	req, _ := http.NewRequest("POST", IMBASEURL+IMMESSAGEPATH, data)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Authorization", auth)
	client := http.DefaultClient
	resp, err := client.Do(req)
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
	if msgResp.Msg_id != "" {
		return nil
	} else {
		return fmt.Errorf("Get message:%s", msgResp.Message)
	}

}
