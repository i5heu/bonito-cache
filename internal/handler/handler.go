package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gotd/contrib/http_range"
	"github.com/i5heu/bonito-cache/internal/config"
	"github.com/i5heu/bonito-cache/internal/helper"
	"github.com/i5heu/bonito-cache/internal/log"
	"github.com/i5heu/bonito-cache/internal/ramCache"
	"github.com/i5heu/bonito-cache/internal/storageCache"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	Conf             config.Config
	DataStoreRAM     *ramCache.DataStore
	DataStoreStorage *storageCache.DataStore
	Log              log.Logger
}

func (h *Handler) HandlerFunc(ctx *fasthttp.RequestCtx) {
	size := uint(0)
	cached := false
	start := time.Now()

	defer func() {
		h.Log.LogRequest(start, string(ctx.RequestURI()), ctx.Response.StatusCode(), cached, size)
	}()

	ctx.Response.Header.Set("Access-Control-Allow-Origin", h.Conf.CORSDomain)
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, HEAD")
	ctx.Response.Header.Set("Cache-Control", "max-age=31536000")
	ctx.Response.Header.Set("Accept-Ranges", "bytes")

	dataResult := h.getData(ctx)
	if dataResult.Error != nil {
		ctx.Response.SetStatusCode(500)
		ctx.Response.Header.Set("Cache-Control", "max-age=0")
		ctx.Response.SetBodyString(dataResult.Error.Error())
		return
	}
	size = dataResult.Size
	cached = dataResult.Cached
	if cached {
		ctx.Response.Header.Set("X-Cached", "true")
	} else {
		ctx.Response.Header.Set("X-Cached", "false")
	}

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
		ctx.Response.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", ranges[0].Start, ranges[0].Start+ranges[0].Length-1, size))
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

	url := config.GetCompleteURL(h.Conf, string(ctx.Path()))
	cachedData, mime := h.DataStoreRAM.GetCacheData(url)
	if cachedData != nil {
		return DataStoreResult{
			Data:   cachedData,
			MIME:   mime,
			Size:   uint(len(cachedData)),
			Cached: true,
			Error:  nil,
		}
	}

	cachedData, mime = h.DataStoreStorage.GetCacheData(url)
	if cachedData != nil {
		return DataStoreResult{
			Data:   cachedData,
			MIME:   mime,
			Size:   uint(len(cachedData)),
			Cached: true,
			Error:  nil,
		}
	}

	newData, sanitizedMime, err := h.getDataFromBackend(url)
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

	// store the requested file in RAM and files storage
	h.DataStoreRAM.CacheData(url, newData, sanitizedMime)
	h.DataStoreStorage.CacheData(url, newData, sanitizedMime)

	return DataStoreResult{
		Data:   newData,
		MIME:   sanitizedMime,
		Size:   uint(len(newData)),
		Cached: false,
		Error:  nil,
	}
}

func (h *Handler) getDataFromBackend(url string) ([]byte, string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}

	sanitizedMime := helper.SanitizeMimeType(res.Header.Get("Content-Type"))

	return bytes, sanitizedMime, nil
}
