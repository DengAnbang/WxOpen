package api

import (
	"net/http"

	"gitee.com/DengAnbang/WxOpen/src/code"
	"gitee.com/DengAnbang/goutils/update"
)

func Run(port string, mux *http.ServeMux) {
	mux.HandleFunc("/", Authentication)
	mux.HandleFunc("/refreshAccessToken", RefreshAccessToken)
	mux.HandleFunc("/menuCreate", MenuCreate)
	//mux.Handle(code.RootName, http.StripPrefix(code.RootName, FileHandler{}))
	mux.Handle("/log/", http.StripPrefix("/log/", http.FileServer(http.Dir(code.LogRootPath))))
	update.UpgradeService(":"+port, mux)
}

type AppHandleFuncErr func(w http.ResponseWriter, r *http.Request) error
