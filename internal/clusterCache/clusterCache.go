package clusterCache

import (
	"errors"

	"github.com/i5heu/bonito-cache/internal/log"

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
	Conf config.Config
	Log  log.Logger
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

func (c *Cluster) ClusterManager() {
	err := c.validateConfig()
	if err != nil {
		panic("Cluster config is invalid")
	}
}
