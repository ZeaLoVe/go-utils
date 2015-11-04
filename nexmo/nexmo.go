package nexmo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	NexmoURL = "https://api.nexmo.com"
	VoiceURL = "/tts/json"
)

type Nexmo struct {
	Url        string
	api_key    string
	api_secret string

	Lg          string
	VoiceGender string
	Repeat      string
	From        string
	To          string
	VoiceMsg    string
	resp        NexmoResponse
}

type NexmoResponse struct {
	Status     int    `json:"status"`
	Error_text string `json:"error_text"`
}

func ToString(words ...string) string {
	res := bytes.Buffer{}
	for _, word := range words {
		res.WriteString(word + "+")
	}
	res.WriteString("over")
	return res.String()
}

func (this *Nexmo) SetKeyAndSecret(key, secret string) {
	this.api_key = key
	this.api_secret = secret
}

func (this *Nexmo) SetTo(to string) {
	this.To = to
}

func (this *Nexmo) SetLanguage(lg string) {
	this.Lg = lg
}

func (this *Nexmo) SetGender(gender string) {
	this.VoiceGender = gender
}

func (this *Nexmo) SetRepeat(times string) {
	this.Repeat = times
}

func (this *Nexmo) SetVoiceMsg(msg string) {
	this.VoiceMsg = msg
}

func (this *Nexmo) Call() error {
	if this.api_key == "" || this.api_secret == "" || this.To == "" || this.VoiceMsg == "" {
		return fmt.Errorf("something needed is not set")
	}
	params := url.Values{}
	reqURL, err := url.Parse(NexmoURL + VoiceURL)

	params.Add("api_key", this.api_key)
	params.Add("api_secret", this.api_secret)
	params.Add("to", this.To)
	params.Add("text", this.VoiceMsg)

	//see https://docs.nexmo.com/api-ref/voice-api/supported-languages  for details
	//zh-cn for Chinese
	if this.Lg != "" {
		params.Add("lg", this.Lg)
	}
	if this.VoiceGender != "" {
		params.Add("voice", this.VoiceGender)
	}
	if this.Repeat != "" {
		params.Add("repeat", this.Repeat)
	}

	reqURL.RawQuery = params.Encode()

	fmt.Println(reqURL.String())
	resp, err := http.Get(reqURL.String())

	if err != nil {
		return fmt.Errorf("Get request error with %s", err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("Get response error with %s", err.Error())
	}

	fmt.Println(string(body))

	json.Unmarshal(body, &this.resp)

	if this.resp.Status != 0 {
		return fmt.Errorf("Response status is %d, error_text:%s", this.resp.Status, this.resp.Error_text)
	} else {
		return nil
	}

}
