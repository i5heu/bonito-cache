package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"simple-S3-cache/config"
	"simple-S3-cache/log"
	"simple-S3-cache/ramCache"

	"github.com/valyala/fasthttp"
)

type Handler struct {
	conf      config.Config
	dataStore *ramCache.DataStore
	log       log.Logger
}

func main() {

	conf := config.GetValues()
	logs := log.New(conf)
	dataStore := ramCache.DataStore{
		Conf: conf,
		Ch:   make(chan ramCache.File, 10000),
	}
	go dataStore.RamFileManager()

	h := Handler{conf: conf, dataStore: &dataStore, log: logs}
	fasthttp.ListenAndServe(":8084", h.handler)
}

func (h *Handler) handler(ctx *fasthttp.RequestCtx) {
	size := uint(0)
	cached := false
	start := time.Now()

	defer func() {
		h.log.LogRequest(start, string(ctx.RequestURI()), ctx.Response.StatusCode(), cached, size)
	}()

	ctx.Response.Header.Set("Access-Control-Allow-Origin", h.conf.CORSDomain)
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET")
	ctx.Response.Header.Set("Cache-Control", "max-age=31536000")

	url := config.GetCompleteURL(h.conf, string(ctx.Path()))
	cachedData := h.dataStore.GetCacheData(url)
	if cachedData != nil {
		cached = true
		size = uint(len(cachedData))
		ctx.Response.SetBody(cachedData)
		return
	}

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		ctx.Response.SetStatusCode(500)
		return
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading response body: ", err)
	}

	size = uint(len(bytes))
	h.dataStore.CacheData(url, bytes)
	ctx.Response.SetBody(bytes)
}
