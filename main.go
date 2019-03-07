package main

import (
	"net/http"

	"gitee.com/DengAnbang/WxOpen/src/api"
	"gitee.com/DengAnbang/WxOpen/src/code"
	"gitee.com/DengAnbang/goutils/loge"
)

//139a0c19f0f90e5a
func main() {
	loge.SetLogPath(code.LogRootPath)
	loge.IsDebug = true
	loge.W("开始服务..")
	mux := http.NewServeMux()
	api.Run("80", mux)
}
