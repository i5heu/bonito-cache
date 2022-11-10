package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/types"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/i5heu/simple-S3-cache/config"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	conf      config.Config
	dataStore *DataStore
}

type DataStore struct {
	// sha256 hash will be divided into 8 uneven parts so we can write while reading almost perfectly
	Data map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string][]byte
}

func main() {
	conf := config.GetValues()
	dataStore := createRamDataStore()
	dataStore.Set("a24f2209335a321e15a1bd72455fec78bd2e87e6a2a4975fca7c4ec7475a4d9d", []byte("test"))

	h := Handler{conf: conf, dataStore: &dataStore}

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

	h.dataStore.Set(sha256Hash, []byte("test"))
	fmt.Println(h.dataStore.Get(sha256Hash))

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

func createRamDataStore() DataStore {
	return DataStore{Data: make(map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string][]byte)}
}

func (d *DataStore) Get(hash string) []byte {
	return d.Data[hash[:1]][hash[1:2]][hash[2:3]][hash[3:4]][hash[4:6]][hash[6:12]][hash[12:24]][hash[24:63]]
}

type dada struct {
	a types.Type
}

func (d *DataStore) Set(hash string, data []byte) {
	createMap(0, d.Data, hash)
	d.Data[hash[:1]][hash[1:2]][hash[2:3]][hash[3:4]][hash[4:6]][hash[6:12]][hash[12:24]][hash[24:63]] = data
}

func createMap(depth int, data interface{}, hash string) interface{} {
	parts := []int{1, 2, 3, 4, 6, 12, 24, 63}
	mapL := data.(map[string]interface{})

	if mapL[hash[parts[depth]:parts[depth+1]]] == nil && depth >= 7 {
		switch depth {
		case 0:
			return createMap(depth+1, make(map[string]map[string]map[string]map[string]map[string]map[string]map[string][]byte), hash)
		case 1:
			return createMap(depth+1, make(map[string]map[string]map[string]map[string]map[string]map[string][]byte), hash)
		case 2:
			return createMap(depth+1, make(map[string]map[string]map[string]map[string]map[string][]byte), hash)
		case 3:
			return createMap(depth+1, make(map[string]map[string]map[string]map[string][]byte), hash)
		case 4:
			return createMap(depth+1, make(map[string]map[string]map[string][]byte), hash)
		case 5:
			return createMap(depth+1, make(map[string]map[string][]byte), hash)
		case 6:
			return createMap(depth+1, make(map[string][]byte), hash)
		default:
			panic("not implemented")
		}

	}

	return data
}
