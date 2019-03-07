package api

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"gitee.com/DengAnbang/WxOpen/src/code"
	"gitee.com/DengAnbang/WxOpen/src/wx"
	"gitee.com/DengAnbang/WxOpen/src/wx/xmlutil"
	"gitee.com/DengAnbang/goutils/httpUtils"
	"gitee.com/DengAnbang/goutils/loge"
)

//手动刷新AccessToken的接口，也会每2小时自动刷新一次
func RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	err := wx.AccessTokenBean.Refresh()
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, wx.AccessTokenBean)
}
func MenuCreate(w http.ResponseWriter, r *http.Request) {
	var menus = make([]wx.Menu, 0)
	menus = append(menus,
		wx.Menu{
			Type: "click",
			Name: "一个按钮",
			Key:  code.KEY_TEST_BUTTON,
		},
		wx.Menu{
			Name: "菜单",
			SubButton: []wx.SubButton{
				{
					Type: "view",
					Name: "搜索1",
					Url:  "http://www.soso.com/",
				},
				{
					Type:     "miniprogram",
					Name:     "小程序",
					Url:      "http://mp.weixin.qq.com",
					Appid:    "wx286b93c14bbf93aa",
					Pagepath: "pages/lunar/index",
				},
				{
					Type: "click",
					Name: "我的资料",
					Key:  code.KEY_MATERIAL,
				},
				{
					Type: "click",
					Name: "生成二维码",
					Key:  code.KEY_CREATE_QR,
				},
			},
		},
		wx.Menu{
			Type: "click",
			Name: "今日歌曲",
			Key:  "V1001_TODAY_MUSIC",
		},
	)
	button := wx.Button{Menu: menus}
	loge.W(wx.MenuCreate(button))

}

//处理微信的认证和收发消息的接口，如果是消息，会分发出去，单独处理
func Authentication(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		signature := httpUtils.GetValueFormRequest(r, "signature")
		timestamp := httpUtils.GetValueFormRequest(r, "timestamp")
		nonce := httpUtils.GetValueFormRequest(r, "nonce")
		echostr := httpUtils.GetValueFormRequest(r, "echostr")
		slice := sort.StringSlice{timestamp, nonce, code.Token}
		sort.Strings(slice)
		str := ""
		for _, s := range slice {
			str += s
		}
		hash := sha1.New()
		hash.Write([]byte(str))
		sprintf := fmt.Sprintf("%x", hash.Sum(nil))
		if sprintf == signature {
			w.Write([]byte(echostr))
			return
		}
		w.Write([]byte(""))
	} else if r.Method == "POST" {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			loge.W(err)
			w.Write([]byte(""))
			return
		}
		loge.WD(string(bytes))
		var messageBean xmlutil.StringMap
		err = xml.Unmarshal(bytes, &messageBean)
		if err != nil {
			loge.W(err)
			w.Write([]byte(""))
			return
		}
		wx.Dispense(w, messageBean)
	}

}
