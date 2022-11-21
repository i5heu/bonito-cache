package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/i5heu/simple-S3-cache/internal/config"
	"github.com/i5heu/simple-S3-cache/internal/helper"
	"github.com/i5heu/simple-S3-cache/internal/log"
	"github.com/i5heu/simple-S3-cache/internal/ramCache"
	"github.com/i5heu/simple-S3-cache/internal/storageCache"

	"github.com/gotd/contrib/http_range"
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
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, HEAD")
	ctx.Response.Header.Set("Cache-Control", "max-age=31536000")
	ctx.Response.Header.Set("Accept-Ranges", "bytes")

	dataResult := h.getData(ctx)
	if dataResult.Error != nil {
		ctx.Response.SetStatusCode(500)
		ctx.Response.SetBodyString(dataResult.Error.Error())
		return
	}
	size = dataResult.Size
	cached = dataResult.Cached

	if ctx.Request.Header.Peek("Range") != nil {
		// parse the range header
		ranges, err := http_range.ParseRange(string(ctx.Request.Header.Peek("Range")), int64(size))
		if err != nil {
			ctx.Response.SetStatusCode(416)
			return
		}
		// we only support one range
		if len(ranges) > 1 {
			ctx.Response.SetStatusCode(416)
			return
		}
		// HTTP 206 (Partial Content)
		ctx.Response.SetStatusCode(206)
		ctx.Response.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", ranges[0].Start, ranges[0].Start+ranges[0].Length, size))
		ctx.Response.Header.Set("Content-Length", strconv.FormatUint(uint64(ranges[0].Length), 10))
		ctx.Response.Header.Set("Content-Type", dataResult.MIME)
	} else {
		ctx.Response.Header.Set("Content-Length", strconv.FormatUint(uint64(size), 10))
		ctx.Response.Header.Set("Content-Type", dataResult.MIME)
	}

	if ctx.IsGet() {
		if ctx.Request.Header.Peek("Range") != nil {
			// parse the range header
			ranges, err := http_range.ParseRange(string(ctx.Request.Header.Peek("Range")), int64(size))
			if err != nil {
				ctx.Response.SetStatusCode(416)
				return
			}

			ctx.Response.SetBody(dataResult.Data[ranges[0].Start : ranges[0].Start+ranges[0].Length])
		} else {
			ctx.Response.SetBody(dataResult.Data)
		}

	}
}

// size and error
type DataStoreResult struct {
	Data   []byte
	MIME   string
	Size   uint
	Cached bool
	Error  error
}

func (h *Handler) getData(ctx *fasthttp.RequestCtx) DataStoreResult {

	url := config.GetCompleteURL(h.conf, string(ctx.Path()))
	cachedData, mime := h.dataStoreRAM.GetCacheData(url)
	if cachedData != nil {
		return DataStoreResult{
			Data:   cachedData,
			MIME:   mime,
			Size:   uint(len(cachedData)),
			Cached: true,
			Error:  nil,
		}
	}

	cachedData, mime = h.dataStoreStorage.GetCacheData(url)
	if cachedData != nil {
		return DataStoreResult{
			Data:   cachedData,
			MIME:   mime,
			Size:   uint(len(cachedData)),
			Cached: true,
			Error:  nil,
		}
	}

	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		ctx.Response.SetStatusCode(500)
		return DataStoreResult{
			Data:   nil,
			MIME:   "",
			Size:   0,
			Cached: false,
			Error:  err,
		}
	}

	if res.StatusCode != 200 {
		ctx.Response.SetStatusCode(res.StatusCode)
		return DataStoreResult{
			Data:   nil,
			MIME:   "",
			Size:   0,
			Cached: false,
			Error:  err,
		}
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading response body: ", err)
	}

	sanitizedMime := helper.SanitizeMimeType(res.Header.Get("Content-Type"))
	h.dataStoreRAM.CacheData(url, bytes, sanitizedMime)
	h.dataStoreStorage.CacheData(url, bytes, sanitizedMime)

	return DataStoreResult{
		Data:   bytes,
		MIME:   sanitizedMime,
		Size:   uint(len(bytes)),
		Cached: false,
		Error:  nil,
	}
}
