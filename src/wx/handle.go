package wx

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"gitee.com/DengAnbang/WxOpen/src/code"
	"gitee.com/DengAnbang/WxOpen/src/wx/xmlutil"
	"gitee.com/DengAnbang/goutils/loge"
)

func Dispense(w http.ResponseWriter, m xmlutil.StringMap) {
	switch m["MsgType"] {
	case "image":
		SendImageMessage(w, m, m["MediaId"])
	case "text":
		SendMessage(w, m, fmt.Sprintf("我收到了消息：%s", m["Content"]))
	case "event":
		if m["Event"] == "CLICK" {
			clickDispense(w, m)
		} else if m["Event"] == "SCAN" {
			scanDispense(w, m)
		}

	}
}
func scanDispense(w http.ResponseWriter, m xmlutil.StringMap) {
	mapAny, err := GetUserMessage(m["EventKey"])
	if err != nil {
		loge.W(err)
		w.Write([]byte(""))
		return
	}
	loge.WD(mapAny)
	SendMessage(w, m, "扫描了"+fmt.Sprint(mapAny["nickname"])+"生成的二维码")
}
func clickDispense(w http.ResponseWriter, m xmlutil.StringMap) {
	switch m["EventKey"] {
	case code.KEY_TEST_BUTTON:
		SendMessage(w, m, fmt.Sprintf("你点击了：%s 按钮", "一个按钮"))
	case code.KEY_MATERIAL:
		mapAny, err := GetUserMessage(m["FromUserName"])
		if err != nil {
			loge.W(err)
			w.Write([]byte(""))
			return
		}
		loge.WD(mapAny)
		SendMessage(w, m, fmt.Sprint(mapAny["nickname"]))
	case code.KEY_SCAN:
		mapAny, err := GetUserMessage(m["EventKey"])
		if err != nil {
			loge.W(err)
			w.Write([]byte(""))
			return
		}
		loge.WD(mapAny)
		SendMessage(w, m, "扫描了"+fmt.Sprint(mapAny["nickname"])+"生成的二维码")
	case code.KEY_CREATE_QR:
		qr, err := CreateQR(m["FromUserName"])
		if err != nil {
			loge.W(err)
			w.Write([]byte(""))
			return
		}
		s := "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + fmt.Sprint(qr["ticket"])
		request, err := http.NewRequest("GET", s, nil)
		resp, err := http.DefaultClient.Do(request)
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			loge.W(err)
			w.Write([]byte(""))
			return
		}
		upload, err := UploaImage(b)
		if err != nil {
			loge.W(err)
			w.Write([]byte(""))
			return
		}
		SendImageMessage(w, m, fmt.Sprint(upload["media_id"]))
	}
}
