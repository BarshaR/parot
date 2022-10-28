package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
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

	server := http.Server{Addr: ":8080", Handler: proxy}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			// handle err
		}
	}()
	fmt.Printf("Proxy running on: %s\n", server.Addr)

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	fmt.Println("Shutting down Parot")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		// handle err
	}
	parotCtx.PrintSummary()

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

func (proxyCtx ParotProxyContext) PrintSummary() {
	totalRecTime := (time.Now().UnixMilli() - proxyCtx.startTime) / 1000
	fmt.Printf("Recorded %d requests over %d seconds\n", proxyCtx.requestHandled, totalRecTime)
}

type LoggingReqHandler struct{}

func (reqHandler LoggingReqHandler) handleRequest(messageNum int, time int64, req []byte) {
	fmt.Printf("Msg # %d Time %d\n%s\n\n", messageNum, time, string(req))
}
