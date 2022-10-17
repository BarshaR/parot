package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			res, err := httputil.DumpRequest(r, true)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(string(res))
			return r, nil
		})
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
