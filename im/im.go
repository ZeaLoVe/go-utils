package im

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

func getToken(acount *Acount) (*UCToken, error) {
	data, err := json.Marshal(&acount)
	if err != nil {
		return nil, err
	}
	//	fmt.Println(string(data))
	req, _ := http.NewRequest("POST", "http://im-agent.web.sdp.101.com/v0.2/api/agents/users/tokens", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	client := http.DefaultClient
	resp, err := client.Do(req)

	token := new(UCToken)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &token)
	//	fmt.Println(string(body))
	if err != nil {
		return nil, err
	} else {
		return token, nil
	}
}

func getIMApi(tos []string, msg string) []byte {
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

func SendMsg(token *UCToken, tos []string, msg string) error {
	auth := fmt.Sprintf("MAC id=\"%s\",nonce=\"%s\",mac=\"%s\"", token.AccessToken, token.Nonce, token.Mac)
	data := string(getIMApi(tos, msg))
	fmt.Println(auth)
	fmt.Println(data)
	req, _ := http.NewRequest("POST", "http://im-agent.web.sdp.101.com/v0.2/api/agents/messages", strings.NewReader(data))
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
	fmt.Println(string(body))
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &msgResp)
	fmt.Println(msgResp.Msg_id)
	if msgResp.Msg_id != "" {
		return nil
	} else {
		return fmt.Errorf("Get message:%s", msgResp.Message)
	}

}
