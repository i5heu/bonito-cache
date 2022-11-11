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
	"sync"

	"github.com/i5heu/simple-S3-cache/config"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	conf      config.Config
	dataStore *DataStore
}

type DataStore struct {
	mu sync.Mutex
	// sha256 hash will be divided into 8 uneven parts so we can write while reading almost perfectly
	Data map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string][]byte
}

func main() {
	conf := config.GetValues()
	dataStore := createRamDataStore()

	h := Handler{conf: conf, dataStore: &dataStore}

	fasthttp.ListenAndServe(":8084", h.handler)
}

func (h *Handler) handler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", h.conf.CORSDomain)
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET")

	url := config.GetCompleteURL(h.conf, string(ctx.Path()))
	cachedData := h.dataStore.GetCacheData(url)
	if cachedData != nil {
		ctx.Response.SetBody(cachedData)
		return
	}

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

	h.dataStore.CacheData(url, bytes)
	ctx.Response.SetBody(bytes)
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
	d.createMap(hash)
	d.Data[hash[:1]][hash[1:2]][hash[2:3]][hash[3:4]][hash[4:6]][hash[6:12]][hash[12:24]][hash[24:63]] = data
}

func (d *DataStore) CacheData(url string, data []byte) {
	hashGen := sha256.New()
	hashGen.Write([]byte(url))
	hash := hex.EncodeToString(hashGen.Sum(nil))

	d.createMap(hash)
	d.Data[hash[:1]][hash[1:2]][hash[2:3]][hash[3:4]][hash[4:6]][hash[6:12]][hash[12:24]][hash[24:63]] = data
}

func (d *DataStore) GetCacheData(url string) []byte {
	hashGen := sha256.New()
	hashGen.Write([]byte(url))
	hash := hex.EncodeToString(hashGen.Sum(nil))

	return d.Data[hash[:1]][hash[1:2]][hash[2:3]][hash[3:4]][hash[4:6]][hash[6:12]][hash[12:24]][hash[24:63]]
}

// there must be a better way to do this
func (d *DataStore) createMap(hash string) {
	parts := []int{0, 1, 2, 3, 4, 6, 12, 24, 63}

	if d.Data[hash[parts[0]:parts[1]]] == nil {
		d.Data[hash[parts[0]:parts[1]]] = make(map[string]map[string]map[string]map[string]map[string]map[string]map[string][]byte)
	}

	if d.Data[hash[parts[0]:parts[1]]][hash[parts[1]:parts[2]]] == nil {
		d.Data[hash[parts[0]:parts[1]]][hash[parts[1]:parts[2]]] = make(map[string]map[string]map[string]map[string]map[string]map[string][]byte)
	}

	if d.Data[hash[parts[0]:parts[1]]][hash[parts[1]:parts[2]]][hash[parts[2]:parts[3]]] == nil {
		d.Data[hash[parts[0]:parts[1]]][hash[parts[1]:parts[2]]][hash[parts[2]:parts[3]]] = make(map[string]map[string]map[string]map[string]map[string][]byte)
	}

	if d.Data[hash[parts[0]:parts[1]]][hash[parts[1]:parts[2]]][hash[parts[2]:parts[3]]][hash[parts[3]:parts[4]]] == nil {
		d.Data[hash[parts[0]:parts[1]]][hash[parts[1]:parts[2]]][hash[parts[2]:parts[3]]][hash[parts[3]:parts[4]]] = make(map[string]map[string]map[string]map[string][]byte)
	}

	if d.Data[hash[parts[0]:parts[1]]][hash[parts[1]:parts[2]]][hash[parts[2]:parts[3]]][hash[parts[3]:parts[4]]][hash[parts[4]:parts[5]]] == nil {
		d.Data[hash[parts[0]:parts[1]]][hash[parts[1]:parts[2]]][hash[parts[2]:parts[3]]][hash[parts[3]:parts[4]]][hash[parts[4]:parts[5]]] = make(map[string]map[string]map[string][]byte)
	}

	if d.Data[hash[parts[0]:parts[1]]][hash[parts[1]:parts[2]]][hash[parts[2]:parts[3]]][hash[parts[3]:parts[4]]][hash[parts[4]:parts[5]]][hash[parts[5]:parts[6]]] == nil {
		d.Data[hash[parts[0]:parts[1]]][hash[parts[1]:parts[2]]][hash[parts[2]:parts[3]]][hash[parts[3]:parts[4]]][hash[parts[4]:parts[5]]][hash[parts[5]:parts[6]]] = make(map[string]map[string][]byte)
	}

	if d.Data[hash[parts[0]:parts[1]]][hash[parts[1]:parts[2]]][hash[parts[2]:parts[3]]][hash[parts[3]:parts[4]]][hash[parts[4]:parts[5]]][hash[parts[5]:parts[6]]][hash[parts[6]:parts[7]]] == nil {
		d.Data[hash[parts[0]:parts[1]]][hash[parts[1]:parts[2]]][hash[parts[2]:parts[3]]][hash[parts[3]:parts[4]]][hash[parts[4]:parts[5]]][hash[parts[5]:parts[6]]][hash[parts[6]:parts[7]]] = make(map[string][]byte)
	}
}
