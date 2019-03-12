package wx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"bytes"

	"gitee.com/DengAnbang/goutils/loge"
)

/**
获取用户的信息
*/
func GetUserMessage(userName string) (map[string]interface{}, error) {
	mapAny := make(map[string]interface{})
	resp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN", AccessTokenBean.AccessToken, userName))
	if err != nil {
		loge.W(err)
		return mapAny, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return mapAny, err
	}
	err = json.Unmarshal(body, &mapAny)
	return mapAny, err
}

/**
创建公众号的菜单
*/
func MenuCreate(button Button) error {
	_, err := POST("https://api.weixin.qq.com/cgi-bin/menu/create?access_token="+AccessTokenBean.AccessToken, button)
	if err != nil {
		return err
	}
	return nil
}
func POST(url string, data interface{}) ([]byte, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(dataBytes)
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		loge.W(err)
		return nil, err
	}
	resp, err := http.DefaultClient.Do(request)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		loge.W(err)
		return nil, err
	}
	return b, nil
}

/**
创建二维码
*/
func CreateQR(fromUserName string) (map[string]interface{}, error) {
	mapQr := make(map[string]interface{})
	s := `{"expire_seconds": 604800, "action_name": "QR_STR_SCENE", "action_info": {"scene": {"scene_str":"` + fromUserName + `"}}}`
	body := bytes.NewReader([]byte(s))
	request, err := http.NewRequest("POST", "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token="+AccessTokenBean.AccessToken, body)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := http.DefaultClient.Do(request)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		loge.W(err)
		return mapQr, err
	}

	err = json.Unmarshal(b, &mapQr)
	if err != nil {
		loge.W(err)
		return mapQr, err
	}
	loge.W(string(b))
	return mapQr, nil
}
