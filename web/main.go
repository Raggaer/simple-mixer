package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"
)

type httpHandlerFunc func(*template.Template, http.ResponseWriter, *http.Request) error

func main() {
	// Load templates
	tpl, err := template.New("SimpleMixer").ParseGlob("views/*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load templates: %v\n", err)
		return
	}

	http.HandleFunc("/", errorHandler(tpl, showHomepage))

	// Create custom server with settings
	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           nil,
		ReadTimeout:       2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       100 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	fmt.Println("Starting HTTP server in :8080")

	if err := httpServer.ListenAndServe(); err != nil {
		fmt.Printf("Unable to start HTTP server: %v\r\n", err)
	}
}

func errorHandler(tpl *template.Template, f httpHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := f(tpl, w, req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
}
