package main

import (
	"net/http"
)

type showHomepageData struct {
	Addr string
	Abi  string
	Rpc  string
}

func showHomepage(ctx *controllerContext) error {
	ctx.res.WriteHeader(http.StatusOK)
	ctx.tpl.ExecuteTemplate(ctx.res, "homepage.html", showHomepageData{
		Rpc:  "http://localhost:8545",
		Abi:  ctx.abi,
		Addr: "0x087F95CccF11F7761Bbd66097e72f730F618Ada2",
	})
	return nil
}
