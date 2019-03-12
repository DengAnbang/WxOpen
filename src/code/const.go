package code

import (
	"gitee.com/DengAnbang/goutils/fileUtil"
)

const (
	OK        = 0  //成功
	NormalErr = -1 //普通错误
)

var (
	CurrentPath = fileUtil.GetCurrentPath()
	RootName    = "/res/"
	RootPath    = CurrentPath + RootName
	LogRootPath = RootPath + "log/"
)

const (
	KEY_TEST_BUTTON  = "KEY_TEST_BUTTON"  //测试按钮
	KEY_MATERIAL     = "KEY_MATERIAL"     //个人资料
	KEY_SEND_NEWS    = "KEY_SEND_NEWS"    //发送图文信息
	KEY_SEND_ARTICLE = "KEY_SEND_ARTICLE" //发送文章信息
	KEY_CREATE_QR    = "KEY_CREATE_QR"    //生成二维码
	KEY_SCAN         = "KEY_SCAN"         //扫描二维码
)
const (
	Token  = "token123456"
	AppID  = "wx9074db6d3187b9ff"
	Secret = "be37a4af66f129bfc4710858e35eae25"
	//Secret = "f476d63c322d31c0839eeb4913de135c"
)
