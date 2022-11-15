package ramCache

import (
	"crypto/sha256"
	"encoding/hex"
	"simple-S3-cache/log"
	"sync"

	"simple-S3-cache/config"
)

type DataStore struct {
	Data sync.Map
	Conf config.Config
	Ch   chan File
	Log  log.Logger
}

func (d *DataStore) Get(hash string) []byte {
	val, ok := d.Data.Load(hash)
	if ok {
		d.Ch <- File{Hash: hash}
		return val.([]byte)
	}

	return nil
}

func (d *DataStore) Set(hash string, data []byte) {
	d.Data.Store(hash, data)
	d.Ch <- File{Hash: hash, Size: uint(len(data))}
}

func (d *DataStore) CacheData(url string, data []byte) {
	hashGen := sha256.New()
	hashGen.Write([]byte(url))
	hash := hex.EncodeToString(hashGen.Sum(nil))

	d.Set(hash, data)
}

func (d *DataStore) GetCacheData(url string) []byte {
	hashGen := sha256.New()
	hashGen.Write([]byte(url))
	hash := hex.EncodeToString(hashGen.Sum(nil))

	return d.Get(hash)
}
