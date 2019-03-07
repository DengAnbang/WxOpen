package wx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

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

//func Upload(mediaType string, body []byte) (map[string]interface{}, error) {
//	err := postFile("adad.jpg", "https://api.weixin.qq.com/cgi-bin/media/upload?access_token="+AccessTokenBean.AccessToken+"&type="+mediaType, body)
//
//	mapQr := make(map[string]interface{})
//	//resopne, err := http.Post("https://api.weixin.qq.com/cgi-bin/media/upload?access_token="+AccessTokenBean.AccessToken+"&type="+mediaType, "multipart/form-data", bytes.NewReader(body))
//	//if err != nil {
//	//	loge.W(err)
//	//	return mapQr, err
//	//}
//	//defer func() {
//	//	resopne.Body.Close()
//	//	loge.W(err)
//	//}()
//	//b, err := ioutil.ReadAll(resopne.Body)
//	//if err != nil {
//	//	loge.W(err)
//	//	return mapQr, err
//	//}
//	//
//	//err = json.Unmarshal(b, &mapQr)
//	//if err != nil {
//	//	loge.W(err)
//	//	return mapQr, err
//	//}
//	return mapQr, err
//}
func UploaImage(body1 []byte) (mediaid map[string]interface{}, err error) {
	mapQr := make(map[string]interface{})

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "qr.jpg")
	if err != nil {
		return
	}
	_, err = io.Copy(part, bytes.NewReader(body1))
	err = writer.Close()
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://qyapi.weixin.qq.com/cgi-bin/media/upload", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	urlQuery := req.URL.Query()
	if err != nil {
		return
	}
	urlQuery.Add("access_token", AccessTokenBean.AccessToken)
	urlQuery.Add("type", "image")

	req.URL.RawQuery = urlQuery.Encode()
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	jsonbody, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(jsonbody, &mapQr)
	if err != nil {
		loge.W(err)
		return mapQr, err
	}
	return
}
func postFile(filename string, targetUrl string, body []byte) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	//iocopy
	_, err = io.Copy(fileWriter, bytes.NewReader(body))
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil
}
