package dingrobot

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const defaultAPI = "https://oapi.dingtalk.com/robot/send"

type Robot struct {
	token   string
	at      robotAt
	keyword string
	secret  string
}

type robotAt struct {
	AtMobiles []string `json:"atMobiles"`
	AtAll     bool     `json:"isAtAll"`
}

type robotRequest struct {
	Type string `json:"msgtype"`
	Text struct {
		Content string `json:"content"`
	} `json:"text,omitempty"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown,omitempty"`
	Link struct {
		Title      string `json:"title"`
		Text       string `json:"text"`
		PictureURL string `json:"picUrl"` // optional
		MessageURL string `json:"messageUrl"`
	} `json:"link,omitempty"`
	At robotAt `json:"at,omitempty"`
}

type robotResponse struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`
}

func New(token string) *Robot {
	return &Robot{
		token: token,
	}
}

func (r *Robot) At(at string) *Robot {
	r.at.AtMobiles = append(r.at.AtMobiles, at)
	r.at.AtAll=false
	return r
}

func (r *Robot) SetKeyWord(key string) *Robot {
	r.keyword = key
	return r
}

func (r *Robot) SetSecret(secret string) *Robot {
	r.secret = secret
	return r
}

func (r *Robot) AtAll(ok bool) *Robot {
	newRobot := *r
	newRobot.at.AtAll = ok
	return &newRobot
}

func (r *Robot) AtMobiles(tels ...string) *Robot {
	newRobot := *r
	newRobot.at.AtMobiles = tels
	return &newRobot
}

func (r *Robot) sendUrl() string {
	if r.keyword != "" {
		return defaultAPI + "?access_token=" + r.token
	} else if r.secret != "" {
		timestamp := time.Now().UnixNano() / 1e6
		stringToSign := fmt.Sprintf("%d\n%s", timestamp, r.secret)

		sign := r.hmacSha256(stringToSign)
		url := fmt.Sprintf("%s&timestamp=%d&sign=%s", defaultAPI+"?access_token="+r.token, timestamp, sign)
		return url
	}
	return ""
}

func (r *Robot) hmacSha256(stringToSign string) string {
	h := hmac.New(sha256.New, []byte(r.secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (r *Robot) postData(request robotRequest) error {
	request.At = r.at
	jdata, _ := json.Marshal(request)
	resp, err := http.Post(r.sendUrl(), "application/json", bytes.NewBuffer(jdata))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	robotResp := robotResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&robotResp); err != nil {
		return err
	}
	if robotResp.Code != 0 {
		return errors.New(robotResp.Message)
	}
	return nil
}

//优先使用 keyword
func (r *Robot) Text(content string) error {
	request := robotRequest{
		Type: "text",
	}
	request.Text.Content = r.keyword + " " + content
	return r.postData(request)
}

func (r *Robot) Markdown(title string, text string) error {
	request := robotRequest{Type: "markdown"}
	request.Markdown.Title = r.keyword + title
	request.Markdown.Text = text
	return r.postData(request)
}

func (r *Robot) Link(title, text string, url string, picUrl string) error {
	request := robotRequest{Type: "link"}
	request.Link.Title = title
	request.Link.Text = r.keyword + text
	request.Link.PictureURL = picUrl
	request.Link.MessageURL = url
	return r.postData(request)
}
