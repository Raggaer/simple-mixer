package main

import (
	"html/template"
	"net/http"
)

func showHomepage(tpl *template.Template, w http.ResponseWriter, req *http.Request) error {
	w.WriteHeader(http.StatusOK)
	tpl.ExecuteTemplate(w, "homepage.html", nil)
	return nil
}
