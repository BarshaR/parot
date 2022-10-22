package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/elazarl/goproxy"
)

func main() {
	parotCtx := ParotProxyContext{
		startTime:      time.Now().UnixMilli(),
		requestHandled: 0,
		requestHandler: LoggingReqHandler{},
	}

	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().Do(&parotCtx)
	log.Fatal(http.ListenAndServe(":8080", proxy))
}

func requestInterceptLogger(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	res, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(res))
	return r, nil
}

type ParotRequestHandler interface {
	handleRequest(messageNumber int, timeDeltaMillis int64, messageBody []byte)
}

type ParotProxyContext struct {
	startTime      int64
	requestHandled int
	requestHandler ParotRequestHandler
}

func (proxyCtx *ParotProxyContext) Handle(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	proxyCtx.requestHandled = proxyCtx.requestHandled + 1
	delta := time.Now().UnixMilli() - proxyCtx.startTime
	res, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Fatal(err)
	}
	proxyCtx.requestHandler.handleRequest(proxyCtx.requestHandled, delta, res)

	return req, nil
}

type LoggingReqHandler struct{}

func (reqHandler LoggingReqHandler) handleRequest(messageNum int, time int64, req []byte) {
	fmt.Printf("Msg # %d Time %d\n%s\n\n", messageNum, time, string(req))
}
