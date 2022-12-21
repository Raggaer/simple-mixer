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
		Addr: "0x213C4dFfFD764765d11FbC067b9Ef89853CCb4a3",
	})
	return nil
}
