package wx

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"io"
	"mime/multipart"
	"os"

	"path"

	"gitee.com/DengAnbang/WxOpen/src/code"
	"gitee.com/DengAnbang/goutils/loge"
	"gitee.com/DengAnbang/WxOpen/src/wx/xmlutil"
	"encoding/xml"
)

var AccessTokenBean AccessToken

func init() {
	AccessTokenBean.gc()
}

//GC回收
func (accessToken *AccessToken) gc() {
	err := accessToken.Refresh()
	if err != nil {
		return
	}
	//定时回收
	time.AfterFunc(time.Duration(accessToken.ExpiresIn)*time.Second, func() { accessToken.gc() })
}
func (accessToken *AccessToken) Refresh() error {
	resp, err := http.PostForm("https://api.weixin.qq.com/cgi-bin/token",
		url.Values{"grant_type": {"client_credential"}, "appid": {code.AppID}, "secret": {code.Secret}})
	if err != nil {
		loge.W(err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		loge.W(err)
		return err
	}
	loge.WD(string(body))
	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		loge.W(err)
		return err
	}
	return nil
}

func UploadImage(filePath string, isUrl bool) (map[string]interface{}, error) {
	mapQr := make(map[string]interface{})
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	paramName := path.Base(filePath)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filePath)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	s := "https://api.weixin.qq.com/cgi-bin/media/upload?access_token=" + AccessTokenBean.AccessToken + "&type=image"
	if isUrl {
		s = "https://api.weixin.qq.com/cgi-bin/media/uploadimg?access_token=" + AccessTokenBean.AccessToken
	}
	request, err := http.NewRequest("POST", s, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		loge.W(err)
	}
	err = json.Unmarshal(b, &mapQr)
	return mapQr, err
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