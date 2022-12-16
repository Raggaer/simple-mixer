package main

import (
	"net/http"
)

func showHomepage(ctx *controllerContext) error {
	ctx.res.WriteHeader(http.StatusOK)
	ctx.tpl.ExecuteTemplate(ctx.res, "homepage.html", nil)
	return nil
}
