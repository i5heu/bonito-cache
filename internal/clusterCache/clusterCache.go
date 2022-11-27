package clusterCache

import (
	"errors"
	"sync"
	"time"

	"github.com/i5heu/bonito-cache/internal/log"
	"github.com/valyala/fasthttp"

	"github.com/i5heu/bonito-cache/internal/config"
)

/*
TODO:
- authenticate other nodes in the cluster
- heartbeat to check if other nodes are still alive
- output all nodes in stdout log
- special non cache request
- round robbin over nodes in cluster to find file
*/

type Cluster struct {
	Conf  config.Config
	Log   log.Logger
	nodes sync.Map
}

func (c *Cluster) validateConfig() error {
	if c.Conf.ClusterKey == "" || len(c.Conf.ClusterKey) < 60 {
		return errors.New("ClusterKey is not empty ot not long enough - make sure it is very secure")
	}
	if c.Conf.ClusterSeed[0:8] != "https://" {
		return errors.New("ClusterSeed dose not begin with https://")
	}

	return nil
}

type NodeStatus struct {
	url           string
	lastHeartbeat time.Time
	failedBeats   int
}

func (c *Cluster) ClusterManager() {
	err := c.validateConfig()
	if err != nil {
		panic("Cluster config is invalid")
	}

	for {
		c.heartbeat()
		c.garbageCollector()
		time.Sleep(5 * time.Second)
	}
}

func (c *Cluster) heartbeat() {
	c.nodes.Range(func(keyI, valueI interface{}) bool {
		key := keyI.(string)
		ns := valueI.(NodeStatus)

		if time.Since(ns.lastHeartbeat) > 12*time.Second {
			c.nodes.Delete(key)
		}

		return false
	})
}
func (c *Cluster) garbageCollector() {
	c.nodes.Range(func(keyI, valueI interface{}) bool {
		key := keyI.(string)
		ns := valueI.(NodeStatus)

		if time.Since(ns.lastHeartbeat) > 12*time.Second {
			c.nodes.Delete(key)
		}

		return false
	})
}

func (c *Cluster) ClusterApiHandler(ctx *fasthttp.RequestCtx) {
	// check if key is correct
}

func (c *Cluster) GetCacheData(url string) ([]byte, string) {

	// iterate through nodes with no-cache req

	return []byte{}, ""
}
