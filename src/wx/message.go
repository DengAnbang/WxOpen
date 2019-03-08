package wx

import (
	"encoding/xml"
	"net/http"

	"gitee.com/DengAnbang/WxOpen/src/wx/xmlutil"
	"gitee.com/DengAnbang/goutils/loge"
)

func SendMessage(w http.ResponseWriter, m xmlutil.StringMap, message string) {
	body := TextRequestBody{
		BaseBody: GetReplyBaseBody("text", m),
		Content:  CDATA{Value: message},
	}
	bytes, err := xml.Marshal(body)
	s := string(bytes)
	loge.W(s)
	if err == nil {
		w.Write(bytes)
		return
	}
	loge.W(err)
	w.Write([]byte(""))
}

func SendImageMessage(w http.ResponseWriter, m xmlutil.StringMap, MediaId string) {
	body := ImageResponseBody{
		BaseBody: GetReplyBaseBody("image", m),
		Image:    struct{ MediaId CDATA }{MediaId: CDATA{Value: MediaId}},
	}
	bytes, err := xml.Marshal(body)
	loge.W(string(bytes))
	if err == nil {
		w.Write(bytes)
		return
	}
	loge.W(err)
	w.Write([]byte(""))
}

func SendNewsMessage(w http.ResponseWriter, m xmlutil.StringMap, articlesItem ArticlesItem) {
	body := NewsResponseBody{
		BaseBody:     GetReplyBaseBody("news", m),
		ArticleCount: 1,
		Articles:     Articles{Item: articlesItem},
	}
	bytes, err := xml.Marshal(body)
	s := string(bytes)
	loge.W(s)
	if err == nil {
		w.Write(bytes)
		return
	}
	loge.W(err)
	w.Write([]byte(""))
}
