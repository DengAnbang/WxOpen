package main

import (
	"net/http"

	"gitee.com/DengAnbang/WxOpen/src/api"
	"gitee.com/DengAnbang/WxOpen/src/code"
	"gitee.com/DengAnbang/goutils/loge"
)

//139a0c19f0f90e5a
func main() {

	//resp, err := http.Get(fmt.Sprintf("http://wx.denganbang.cn/"))
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(body))

	loge.SetLogPath(code.LogRootPath)
	loge.IsDebug = true
	loge.W("开始服务..")
	mux := http.NewServeMux()
	api.Run("80", mux)
}
