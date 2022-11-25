package main

import (
	"github.com/i5heu/bonito-cache/internal/clusterCache"
	"github.com/i5heu/bonito-cache/internal/config"
	"github.com/i5heu/bonito-cache/internal/handler"
	"github.com/i5heu/bonito-cache/internal/log"
	"github.com/i5heu/bonito-cache/internal/ramCache"
	"github.com/i5heu/bonito-cache/internal/storageCache"

	"github.com/valyala/fasthttp"
)

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

	if conf.ClusterActive {
		cluster := clusterCache.Cluster{
			Conf: conf,
			Log:  logs,
		}
		go cluster.ClusterManager()
	}

	h := handler.Handler{
		Conf:             conf,
		DataStoreRAM:     &dataStoreRAM,
		DataStoreStorage: &dataStoreStorage,
		Log:              logs,
	}
	fasthttp.ListenAndServe(":8084", h.HandlerFunc)
}
