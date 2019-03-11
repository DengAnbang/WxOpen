package wx

import (
	"encoding/xml"
	"net/http"

	"bytes"
	"encoding/json"
	"io/ioutil"

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

func SendNewsMessage(w http.ResponseWriter, m xmlutil.StringMap, articlesItem NewsArticlesItem) {
	body := NewsResponseBody{
		BaseBody:     GetReplyBaseBody("news", m),
		ArticleCount: 1,
		Articles:     NewsArticles{Item: articlesItem},
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

/**
上传素材
*/
func UploadArticleMessage(w http.ResponseWriter, articles Articles) xmlutil.StringMap {
	stringMap := make(xmlutil.StringMap, 0)
	bytess, err := json.Marshal(articles)
	if err != nil {
		loge.W(err)
		w.Write([]byte(""))
		return stringMap
	}
	//fmt.Fprint(w, string(bytess))
	body := bytes.NewReader(bytess)
	request, err := http.NewRequest("POST", "https://api.weixin.qq.com/cgi-bin/media/uploadnews?access_token="+AccessTokenBean.AccessToken, body)
	if err != nil {
		loge.W(err)
		w.Write([]byte(""))
		return stringMap
	}
	resp, err := http.DefaultClient.Do(request)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	loge.W(string(b))
	if err != nil {
		loge.W(err)
		w.Write([]byte(""))
		return stringMap
	}

	err = xml.Unmarshal(b, &stringMap)
	if err != nil {
		loge.W(err)
		w.Write([]byte(""))
		return stringMap
	}
	return stringMap
}
