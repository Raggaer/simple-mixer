package main

import (
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	contractAddress = "0x0"
)

type httpHandlerFunc func(*controllerContext) error

type controllerContext struct {
	abi    string
	res    http.ResponseWriter
	priv   *privateKey
	req    *http.Request
	tpl    *template.Template
	client *ethclient.Client
}

type applicationFlags struct {
	privateKey      string
	contractAddress string
}

func main() {
	flags := parseApplicationFlags()
	contractAddress = flags.contractAddress

	// Load templates
	tpl, err := template.New("SimpleMixer").ParseGlob("views/*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load templates: %v\n", err)
		return
	}

	// Load private key
	privKey, err := parsePrivateKey(flags.privateKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse private key: %v\n", err)
		return
	}

	// Load contract ABI
	abi, err := loadContractABI(filepath.Join("abi", "SimpleMixer.json"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load contract ABI JSON file: %v\n", err)
		return
	}

	// Connect to RPC server
	client, err := createEthClient("http://localhost:8545")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to RPC server: %v\n", err)
		return
	}

	http.HandleFunc("/", errorHandler(client, privKey, abi, tpl, showHomepage))
	http.HandleFunc("/api/sign", errorHandler(client, privKey, abi, tpl, sendSignature))
	fs := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/static/", staticHandler(http.StripPrefix("/static", fs)))

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

func parseApplicationFlags() *applicationFlags {
	f := applicationFlags{}
	flag.StringVar(&f.privateKey, "privateKey", "0x0", "Private key for message signing")
	flag.StringVar(&f.contractAddress, "contractAddress", "0x0", "Address where the mixer contract is deployed")
	flag.Parse()
	return &f
}

func staticHandler(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fs.ServeHTTP(w, req)
	}
}

// Handle any controller action by returning an error
func errorHandler(client *ethclient.Client, privKey *privateKey, abi string, tpl *template.Template, f httpHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := &controllerContext{
			abi:    abi,
			priv:   privKey,
			client: client,
			tpl:    tpl,
			req:    req,
			res:    w,
		}

		if err := f(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
}

// Connects to the given RPC server
func createEthClient(rpc string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(rpc)
	return client, err
}

// Load ABI on startup so we dont need to load the file on every request
func loadContractABI(p string) (string, error) {
	content, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
