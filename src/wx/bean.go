package wx

import (
	"encoding/xml"
	"time"

	"gitee.com/DengAnbang/WxOpen/src/wx/xmlutil"
)

//认证access_token
type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
}

type Menu struct {
	Name      string      `json:"name,omitempty"`
	Type      string      `json:"type,omitempty"`
	Key       string      `json:"key,omitempty"`
	Url       string      `json:"url,omitempty"`
	Appid     string      `json:"appid,omitempty"`
	Pagepath  string      `json:"pagepath,omitempty"`
	MediaId   string      `json:"media_id,omitempty"`
	SubButton []SubButton `json:"sub_button,omitempty"`
}
type SubButton struct {
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Key      string `json:"key,omitempty"`
	Url      string `json:"url,omitempty"`
	Appid    string `json:"appid,omitempty"`
	Pagepath string `json:"pagepath,omitempty"`
	MediaId  string `json:"media_id,omitempty"`
}
type Button struct {
	Menu []Menu `json:"button"`
}
type CDATA struct {
	Value string `xml:",cdata"`
}

//通用的属性
type BaseBody struct {
	XMLName      xml.Name `xml:"xml"`
	MsgId        int      `xml:"-"`
	CreateTime   time.Duration
	ToUserName   CDATA
	FromUserName CDATA
	MsgType      CDATA
}

func GetReplyBaseBody(msgType string, m xmlutil.StringMap) BaseBody {
	baseBody := BaseBody{
		CreateTime:   time.Duration(time.Now().Unix()),
		FromUserName: CDATA{Value: m["ToUserName"]},
		ToUserName:   CDATA{Value: m["FromUserName"]},
		MsgType:      CDATA{Value: msgType},
	}
	return baseBody
}

type TextRequestBody struct {
	BaseBody
	Content CDATA
}
type ImageResponseBody struct {
	BaseBody
	Image struct {
		MediaId CDATA
	}
}

type NewsResponseBody struct {
	BaseBody
	ArticleCount int32
	Articles     Articles
}

type Articles struct {
	Item ArticlesItem `xml:"item"`
}

type ArticlesItem struct {
	Title       CDATA
	Description CDATA
	PicUrl      CDATA
	Url         CDATA
}
