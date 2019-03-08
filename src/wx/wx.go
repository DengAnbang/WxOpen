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

//func UploaImage(body1 []byte) (mediaid map[string]interface{}, err error) {
//	mapQr := make(map[string]interface{})
//	body := &bytes.Buffer{}
//	writer := multipart.NewWriter(body)
//	part, err := writer.CreateFormFile("image", "qr.jpg")
//	if err != nil {
//		return
//	}
//	_, err = io.Copy(part, bytes.NewReader(body1))
//	err = writer.Close()
//	if err != nil {
//		return
//	}
//
//	req, err := http.NewRequest("POST", "https://qyapi.weixin.qq.com/cgi-bin/media/upload", body)
//	req.Header.Add("Content-Type", writer.FormDataContentType())
//	urlQuery := req.URL.Query()
//	if err != nil {
//		return
//	}
//	urlQuery.Add("access_token", AccessTokenBean.AccessToken)
//	urlQuery.Add("type", "image")
//
//	req.URL.RawQuery = urlQuery.Encode()
//	client := http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		return
//	}
//	defer res.Body.Close()
//	jsonbody, _ := ioutil.ReadAll(res.Body)
//	err = json.Unmarshal(jsonbody, &mapQr)
//	if err != nil {
//		loge.W(err)
//		return mapQr, err
//	}
//	return
//}
//type MediaUpload struct {
//	ErrCode   int    `json:"errcode"`
//	ErrMgs    string `json:"errmsg"`
//	Type      string `json:"type"`
//	MediaID   string `json:"media_id"`
//	CreatedAt string `json:"created_at"`
//}
//
//func UploaImage(token, imagePath string) (mediaid string, err error) {
//	file, err := os.Open(imagePath)
//	if err != nil {
//		return
//	}
//	defer file.Close()
//	body := &bytes.Buffer{}
//	writer := multipart.NewWriter(body)
//	part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
//	if err != nil {
//		return
//	}
//	_, err = io.Copy(part, file)
//	err = writer.Close()
//	if err != nil {
//		return
//	}
//
//	req, err := http.NewRequest("POST", "https://api.weixin.qq.com/cgi-bin/media/upload", body)
//	req.Header.Add("Content-Type", writer.FormDataContentType())
//	urlQuery := req.URL.Query()
//	if err != nil {
//		return
//	}
//	urlQuery.Add("access_token", token)
//	urlQuery.Add("type", "image")
//
//	req.URL.RawQuery = urlQuery.Encode()
//	client := http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		return
//	}
//	defer res.Body.Close()
//	jsonbody, _ := ioutil.ReadAll(res.Body)
//	media := MediaUpload{}
//	err = json.Unmarshal(jsonbody, &media)
//	if err != nil {
//		return
//	}
//	if media.MediaID == "" {
//		err = errors.New(media.ErrMgs)
//	}
//	mediaid = media.MediaID
//	return
//}

func UploaImage(filePath string) (map[string]interface{}, error) {
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
	request, err := http.NewRequest("POST", "https://api.weixin.qq.com/cgi-bin/media/upload?access_token="+AccessTokenBean.AccessToken+"&type=image", body)
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
func newfileUploadRequest(uri string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, path)
	if err != nil {
		return nil, err
	}
	// 这里的io.Copy实现,会把file文件都读取到内存里面，然后当做一个buffer传给NewRequest. 对于大文件来说会占用很多内存
	_, err = io.Copy(part, file)

	//for key, val := range params {
	//	_ = writer.WriteField(key, val)
	//}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", uri, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request, err
}
