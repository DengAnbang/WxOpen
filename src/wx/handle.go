package wx

import (
	"fmt"
	"net/http"

	"io/ioutil"

	"os"

	"gitee.com/DengAnbang/WxOpen/src/code"
	"gitee.com/DengAnbang/WxOpen/src/wx/xmlutil"
	"gitee.com/DengAnbang/goutils/loge"
	"gitee.com/DengAnbang/goutils/utils"
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
	case code.KEY_SEND_NEWS:
		articlesItem := ArticlesItem{Title: CDATA{Value: "这是标题"},
			Description: CDATA{Value: "这是标题的描述..."},
			PicUrl:      CDATA{Value: "http://mmbiz.qpic.cn/mmbiz_jpg/hvWibETJA6ZON5sQMalx0NicA3rwFjbDoJERwJw1qrtDnAYWHqo5rhY6cScib1FXytuzNgZicCHAibMgcaE4ObGT2Bw/0"},
			Url:         CDATA{Value: "https://www.baidu.com/"},
		}
		SendNewsMessage(w, m, articlesItem)
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
		filePath := utils.NewUUID() + ".png"
		err = ioutil.WriteFile(filePath, b, 0666)
		if err != nil {
			loge.W(err)
			w.Write([]byte(""))
			return
		}
		defer os.Remove(filePath)
		mapQr, err := UploaImage(filePath)
		if err != nil {
			loge.W(err)
			w.Write([]byte(""))
			return
		}
		loge.W(mapQr)
		SendImageMessage(w, m, fmt.Sprint(mapQr["media_id"]))
	}
}
