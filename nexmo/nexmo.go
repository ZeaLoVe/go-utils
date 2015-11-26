package nexmo

/**
usage:
	var nexmo Nexmo
	nexmo.SetKeyAndSecret("your_key", "your_secret")
	nexmo.SetVoiceMsg("your-words")
	nexmo.SetTo("86***********") //your-phone-number ,86 for china
	nexmo.SetRepeat("2") //repeat-times
	nexmo.SetLanguage("zh-cn") //default en-us
	//	nexmo.SetGender("male") //default femail

	err := nexmo.Call()
	if err != nil {
		fmt.Println(err.Error())
	}
**/
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

var (
	NexmoURL = "https://api.nexmo.com"
	VoiceURL = "/tts/json"
)

type Nexmo struct {
	sync.Mutex
	Url        string
	api_key    string
	api_secret string

	Lg          string
	VoiceGender string
	Repeat      string
	From        string
	To          string
	VoiceMsg    string
	Status_url  string
	Resp        NexmoResponse
}

type NexmoResponse struct {
	To         string `json:"to"`
	Call_id    string `json:"call_id"`
	Status     int    `json:"status"`
	Error_text string `json:"error_text"`
}

func (this *Nexmo) SetKeyAndSecret(key, secret string) {
	this.api_key = key
	this.api_secret = secret
}

func (this *Nexmo) SetToChinaZoneCode(to string) {
	this.To = "86" + to
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

func (this *Nexmo) SetCallBack(status_url string) {
	this.Status_url = status_url
}

func (this *Nexmo) Call() (error, NexmoResponse) {
	if this.api_key == "" || this.api_secret == "" || this.To == "" || this.VoiceMsg == "" {
		return fmt.Errorf("something need is missing"), NexmoResponse{}
	}
	this.Lock()
	defer this.Unlock()
	params := url.Values{}
	reqURL, err := url.Parse(NexmoURL + VoiceURL)

	params.Add("api_key", this.api_key)
	params.Add("api_secret", this.api_secret)
	params.Add("to", this.To)
	params.Add("text", this.VoiceMsg)
	if this.Status_url != "" {
		params.Add("callback", this.Status_url)
	}

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

	resp, err := http.Get(reqURL.String())

	if err != nil {
		return fmt.Errorf("Get request error with %s", err.Error()), NexmoResponse{}
	}

	body, err := ioutil.ReadAll(resp.Body)

	//	fmt.Println(string(body))

	if err != nil {
		return fmt.Errorf("Get response error with %s", err.Error()), NexmoResponse{}
	}

	json.Unmarshal(body, &this.Resp)

	if this.Resp.Status != 0 {
		return fmt.Errorf("Response status is %d, error_text:%s", this.Resp.Status, this.Resp.Error_text), this.Resp
	} else {
		return nil, this.Resp
	}

}
