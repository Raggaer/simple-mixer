package main

import (
	"net/http"
)

type showHomepageData struct {
	Addr string
	Abi  string
	Rpc  string
}

// Renders the main page
func showHomepage(ctx *controllerContext) error {
	ctx.res.WriteHeader(http.StatusOK)
	ctx.tpl.ExecuteTemplate(ctx.res, "homepage.html", showHomepageData{
		Rpc:  "http://localhost:8545",
		Abi:  ctx.abi,
		Addr: contractAddress,
	})
	return nil
}
