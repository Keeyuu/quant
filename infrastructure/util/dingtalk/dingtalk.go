package dingtalk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"time"
)

type DingtalkRobot struct {
	NotifyUrl string
	SecretKey string
}

func (r *DingtalkRobot) SendMsg(msg *DingtalkMsg) error {
	notifyUrl := r.NotifyUrl
	if r.SecretKey != "" {
		ts, sign := r.calcSign(r.SecretKey)
		notifyUrl = fmt.Sprintf("%v&timestamp=%v&sign=%v", r.NotifyUrl, ts, sign)
	}
	reader := bytes.NewReader(msg.genMsg())
	resp, err := http.Post(notifyUrl, "application/json", reader)
	if err != nil {
		return err
	}

	var resData struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if b, err := ioutil.ReadAll(resp.Body); err != nil {
		return err
	} else if err := json.Unmarshal(b, &resData); err != nil {
		return err
	}

	if resData.ErrCode != 0 {
		return errors.New(resData.ErrMsg)
	}
	return nil
}

func (*DingtalkRobot) calcSign(secretKey string) (ts int64, sign string) {
	ms := math.Round(float64(time.Now().UnixNano()) / 1e6)
	stringToSignEnc := []byte(fmt.Sprintf("%v\n%v", int64(ms), secretKey))
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(stringToSignEnc)
	b := h.Sum(nil)
	encodedStr := base64.StdEncoding.EncodeToString(b)
	escapedStr := url.QueryEscape(encodedStr)
	return int64(ms), escapedStr
}

type DingtalkMsg struct {
	typ string
	msg map[string]interface{}
	at  struct {
		IsAtAll bool `json:"isAtAll, omitempty"`
	}
}

func (l *DingtalkMsg) genMsg() []byte {
	data := make(map[string]interface{})
	data["msgtype"] = l.typ
	data[l.typ] = l.msg

	// at
	atMap := map[string]interface{}{}
	if l.at.IsAtAll {
		atMap["isAtAll"] = true
	}
	data["at"] = atMap
	b, _ := json.Marshal(data)
	return b
}

func NewTextMsg(text string) *DingtalkMsg {
	var obj DingtalkMsg
	obj.typ = "text"
	obj.msg = map[string]interface{}{
		"content": text,
	}
	return &obj
}

func NewLinkMsg(title, text, msgUrl, picUrl string) *DingtalkMsg {
	var obj DingtalkMsg
	obj.typ = "link"
	obj.msg = map[string]interface{}{
		"title":      title,
		"text":       text,
		"messageUrl": msgUrl,
		"picUrl":     picUrl,
	}
	return &obj
}

func NewMarkdownMsg(title, text string) *DingtalkMsg {
	var obj DingtalkMsg
	obj.typ = "markdown"
	obj.msg = map[string]interface{}{
		"title": title,
		"text":  text,
	}
	return &obj
}

type KVPair struct {
	Key   string
	Value interface{}
}

// 对Markdown消息的封装
func NewKVMsg(title string, kvs []KVPair) *DingtalkMsg {
	data := "### " + title
	for _, v := range kvs {
		data += fmt.Sprintf("  \n> **%s**: %v", v.Key, v.Value)
	}
	return NewMarkdownMsg(title, data)
}
