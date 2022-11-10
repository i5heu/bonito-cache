package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/i5heu/simple-S3-cache/config"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	conf config.Config
}

func main() {
	conf := config.GetValues()
	h := Handler{conf: conf}

	fasthttp.ListenAndServe(":8084", h.handler)
}

func (h *Handler) handler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world")

	url := config.GetCompleteURL(h.conf, string(ctx.Path()))

	// get sha256 hash from url
	hash := sha256.New()
	hash.Write([]byte(url))
	sha256Hash := hex.EncodeToString(hash.Sum(nil))
	fmt.Println(sha256Hash)

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	ctx.Response.SetBody(bytes)

	ctx.Response.Header.Set("Access-Control-Allow-Origin", h.conf.CORSDomain)
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET")
}
