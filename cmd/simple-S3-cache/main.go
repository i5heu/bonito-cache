package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/i5heu/simple-S3-cache/internal/config"
	"github.com/i5heu/simple-S3-cache/internal/helper"
	"github.com/i5heu/simple-S3-cache/internal/log"
	"github.com/i5heu/simple-S3-cache/internal/ramCache"
	"github.com/i5heu/simple-S3-cache/internal/storageCache"

	"github.com/valyala/fasthttp"
)

type Handler struct {
	conf             config.Config
	dataStoreRAM     *ramCache.DataStore
	dataStoreStorage *storageCache.DataStore
	log              log.Logger
}

func main() {

	conf := config.GetValues()
	logs := log.New(conf)
	dataStoreRAM := ramCache.DataStore{
		Conf: conf,
		Ch:   make(chan ramCache.File, 10000),
		Log:  logs,
	}
	go dataStoreRAM.RamFileManager()

	dataStoreStorage := storageCache.DataStore{
		Conf: conf,
		Log:  logs,
	}
	go dataStoreStorage.StorageFileManager()

	h := Handler{conf: conf, dataStoreRAM: &dataStoreRAM, dataStoreStorage: &dataStoreStorage, log: logs}
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
	cachedData, mime := h.dataStoreRAM.GetCacheData(url)
	if cachedData != nil {
		cached = true
		size = uint(len(cachedData))
		ctx.Response.Header.Set("Content-Type", mime)
		ctx.Response.SetBody(cachedData)
		return
	}

	cachedData, mime = h.dataStoreStorage.GetCacheData(url)
	if cachedData != nil {
		cached = true
		h.dataStoreRAM.CacheData(url, cachedData, mime)
		size = uint(len(cachedData))
		ctx.Response.Header.Set("Content-Type", mime)
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

	if res.StatusCode != 200 {
		ctx.Response.SetStatusCode(res.StatusCode)
		return
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading response body: ", err)
	}

	size = uint(len(bytes))
	sanitizedMime := helper.SanitizeMimeType(res.Header.Get("Content-Type"))
	h.dataStoreRAM.CacheData(url, bytes, sanitizedMime)
	h.dataStoreStorage.CacheData(url, bytes, sanitizedMime)
	ctx.Response.Header.Set("Content-Type", sanitizedMime)
	ctx.Response.SetBody(bytes)
}
