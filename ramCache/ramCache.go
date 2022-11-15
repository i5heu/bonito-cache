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

type FileStore struct {
	data []byte
	mime string
}

func (d *DataStore) Get(hash string) ([]byte, string) {
	val, ok := d.Data.Load(hash)
	if ok {
		d.Ch <- File{Hash: hash}
		return val.(FileStore).data, val.(FileStore).mime
	}

	return nil, ""
}

func (d *DataStore) Set(hash string, data []byte, mime string) {
	d.Data.Store(hash, FileStore{data: data, mime: mime})
	d.Ch <- File{Hash: hash, Size: uint(len(data)), MIME: mime}
}

func (d *DataStore) CacheData(url string, data []byte, mime string) {
	hashGen := sha256.New()
	hashGen.Write([]byte(url))
	hash := hex.EncodeToString(hashGen.Sum(nil))

	d.Set(hash, data, mime)
}

func (d *DataStore) GetCacheData(url string) ([]byte, string) {
	hashGen := sha256.New()
	hashGen.Write([]byte(url))
	hash := hex.EncodeToString(hashGen.Sum(nil))

	return d.Get(hash)
}
