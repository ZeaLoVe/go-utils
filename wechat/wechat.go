//all code from https://github.com/Yanjunhui/chat.git
package wechat

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	CorpId         string
	AgentId        string
	SecretKey      string
	EncodingAESKey string

	TokenCache *cache.Cache
)

func SetConfig(corpid string, agentid string, secret string, aeskey string) {
	CorpId = corpid
	AgentId = agentid
	SecretKey = secret
	EncodingAESKey = aeskey
}

func init() {
	TokenCache = cache.New(6000*time.Second, 5*time.Second)
}

//发送信息
type Content struct {
	Content string `json:"content"`
}

type MsgPost struct {
	ToUser  string  `json:"touser"`
	MsgType string  `json:"msgtype"`
	AgentID int     `json:"agentid"`
	Text    Content `json:"text"`
}

func SendWxMsg(toUser string, content string) error {
	if userList := strings.Split(toUser, ","); len(userList) > 1 {
		toUser = strings.Join(userList, "|")
	}

	text := Content{}
	text.Content = content

	msg := MsgPost{
		ToUser:  toUser,
		MsgType: "text",
		AgentID: StringToInt(AgentId),
		Text:    text,
	}

	token, found := TokenCache.Get("token")
	if !found {
		return fmt.Errorf("token获取失败!")
	}
	accessToken, ok := token.(AccessToken)
	if !ok {
		return fmt.Errorf("token解析失败!")
	}

	url := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + accessToken.AccessToken

	result, err := WxPost(url, msg)
	if err != nil {
		return fmt.Errorf("请求微信失败: %v", err)
	}
	log.Printf("发送信息给%s, 信息内容: %s, 微信返回结果: %v", toUser, content, result)
	return nil
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

//从微信获取 AccessToken
func GetAccessTokenFromWeixin() {

	for {
		if CorpId == "" || SecretKey == "" {
			log.Printf("corpId或者secret 获取失败, 请检查配置文件")
			return
		}

		WxAccessTokenUrl := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + CorpId + "&corpsecret=" + SecretKey

		tr := &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
		}
		client := &http.Client{Transport: tr}
		result, err := client.Get(WxAccessTokenUrl)
		if err != nil {
			log.Printf("获取微信 Token 返回数据错误: %v", err)
			time.Sleep(time.Minute)
			continue
		}

		res, err := ioutil.ReadAll(result.Body)

		if err != nil {
			log.Printf("获取微信 Token 返回数据错误: %v", err)
			time.Sleep(time.Minute)
			continue
		}
		newAccess := AccessToken{}
		err = json.Unmarshal(res, &newAccess)
		if err != nil {
			log.Printf("获取微信 Token 返回数据解析 Json 错误: %v", err)
			time.Sleep(time.Minute)
			continue
		}

		if newAccess.ExpiresIn == 0 || newAccess.AccessToken == "" {
			log.Printf("获取微信错误代码: %v, 错误信息: %v", newAccess.ErrCode, newAccess.ErrMsg)
			time.Sleep(5 * time.Minute)
		}

		//延迟时间
		TokenCache.Set("token", newAccess, time.Duration(newAccess.ExpiresIn)*time.Second)
		log.Printf("微信 Token 更新成功: %s,有效时间: %v", newAccess.AccessToken, newAccess.ExpiresIn)
		time.Sleep(time.Duration(newAccess.ExpiresIn-100) * time.Second)
	}

}

//微信请求数据
func WxPost(url string, data MsgPost) (string, error) {
	jsonBody, err := encodeJson(data)
	if err != nil {
		return "", err
	}

	r, err := http.Post(url, "application/json;charset=utf-8", bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}

//string 类型转 int
func StringToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("agent 类型转换失败, 请检查配置文件中 agentid 配置是否为纯数字(%v)", err)
		return 0
	}
	return n
}

//json序列化(禁止 html 符号转义)
func encodeJson(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
