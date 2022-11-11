package ramCache

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

type DataStore struct {
	Data sync.Map
}

func Create() DataStore {
	return DataStore{}
}

func (d *DataStore) Get(hash string) []byte {
	val, ok := d.Data.Load(hash)
	if ok {
		return val.([]byte)
	}

	return nil
}

func (d *DataStore) Set(hash string, data []byte) {
	d.Data.Store(hash, data)
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
