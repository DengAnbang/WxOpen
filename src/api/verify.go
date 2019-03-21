package api

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"bytes"
	"encoding/json"

	"gitee.com/DengAnbang/WxOpen/src/code"
	"gitee.com/DengAnbang/WxOpen/src/wx"
	"gitee.com/DengAnbang/WxOpen/src/wx/xmlutil"
	"gitee.com/DengAnbang/goutils/httpUtils"
	"gitee.com/DengAnbang/goutils/loge"
)

func Test1(w http.ResponseWriter, r *http.Request) {
	body := bytes.NewReader([]byte(`{
   "touser":"oLg5g1U2t04ZWVejK9LtfJ_PHDMM",
   "mpnews":{
     "media_id":"u4DTvfet1YK-HS0vIjh7kV5I9jZMNYzOkDIgQ0hh-9B2pCusGOKAd_ZbwhqAik4j"
    },
   "msgtype":"mpnews"
}`))
	request, err := http.NewRequest("POST", "https://api.weixin.qq.com/cgi-bin/message/mass/preview?access_token="+wx.AccessTokenBean.AccessToken, body)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	resp, err := http.DefaultClient.Do(request)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	fmt.Fprint(w, string(b))
}
func Test2(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Test2")
}
func Test(w http.ResponseWriter, r *http.Request) {
	//tv4-O5UjqdhUjWCDCQ8i5cPl9kjYFwR3tkD42FJ7rmNQ2g_ZpWOMzECa2-43yEK6
	image, err := wx.UploadImage(`E:\code\golang\src\gitee.com\DengAnbang\WxOpen\200812308231244_2.jpg`, false)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	article := make([]wx.Article, 0)
	article = append(article, wx.Article{
		ThumbMediaId:       fmt.Sprint(image["media_id"]),
		Author:             "Author",
		Title:              "title",
		Content:            "Content",
		ContentSourceUrl:   "www.baidu.com",
		Digest:             "Digest",
		NeedOpenComment:    1,
		OnlyFansCanComment: 0,
		ShowCoverPic:       1,
	})
	articles := wx.Articles{
		Article: article,
	}

	bytess, err := json.Marshal(articles)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, string(bytess))
	body := bytes.NewReader(bytess)
	request, err := http.NewRequest("POST", "https://api.weixin.qq.com/cgi-bin/media/uploadnews?access_token="+wx.AccessTokenBean.AccessToken, body)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	resp, err := http.DefaultClient.Do(request)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, string(b))

}

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
			Name: "菜单2",
			SubButton: []wx.SubButton{
				{
					Type: "click",
					Name: "发送图文消息",
					Key:  code.KEY_SEND_NEWS,
				},
				{
					Type: "click",
					Name: "发送文章消息",
					Key:  code.KEY_SEND_ARTICLE,
				},
			},
		},
	)
	button := wx.Button{Menu: menus}
	err := wx.MenuCreate(button)
	if err != nil {
		fmt.Fprintf(w, "更新按钮失败：%v", err)
		loge.W(err)
		return
	}
	fmt.Fprint(w, "更新按钮完成")

}

//处理微信的认证和收发消息的接口，如果是消息，会分发出去，单独处理
func Authentication(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		signature := httpUtils.GetValueFormRequest(r, "signature")
		timestamp := httpUtils.GetValueFormRequest(r, "timestamp")
		nonce := httpUtils.GetValueFormRequest(r, "nonce")
		echostr := httpUtils.GetValueFormRequest(r, "echostr")
		slice := sort.StringSlice{timestamp, nonce, code.Token}
		r.ParseForm()
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
		bytess, err := ioutil.ReadAll(r.Body)
		if err != nil {
			loge.W(err)
			w.Write([]byte(""))
			return
		}
		loge.WD(string(bytess))
		var messageBean xmlutil.StringMap
		err = xml.Unmarshal(bytess, &messageBean)
		if err != nil {
			loge.W(err)
			w.Write([]byte(""))
			return
		}
		wx.Dispense(w, messageBean)
	}

}
