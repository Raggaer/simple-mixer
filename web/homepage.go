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
		Addr: "0xc4Ba3D829821F1569F1833D1874734bBf403255e",
	})
	return nil
}
