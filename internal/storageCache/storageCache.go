package storageCache

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/i5heu/simple-S3-cache/internal/config"
	"github.com/i5heu/simple-S3-cache/internal/log"
)

type DataStore struct {
	Conf config.Config
	Log  log.Logger
}

func (d *DataStore) Get(hash string) ([]byte, string) {
	data, err := os.ReadFile(d.GetPath(hash))
	// if no file found, return empty data
	if os.IsNotExist(err) {
		return nil, ""
	}
	if err != nil {
		fmt.Println(err)
		return nil, ""
	}

	dataRaw := data[260:]
	mimeByte := data[:260]
	mime := string(mimeByte[:bytes.Index(mimeByte, []byte{0})])

	err = os.Chtimes(d.GetPath(hash), time.Now(), time.Now())
	if err != nil {
		fmt.Println(err)
	}

	return dataRaw, mime
}

func (d *DataStore) Set(hash string, rawData []byte, mime string) {
	mimeBtytes := append([]byte(mime), make([]byte, 260-len([]byte(mime)))...)
	data := append(mimeBtytes, rawData...)

	err := os.WriteFile(d.GetPath(hash), data, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (d *DataStore) delete(hash string) {
	err := os.Remove(d.GetPath(hash))
	if err != nil {
		fmt.Println(err)
	}
}

func (d *DataStore) GetCacheData(url string) ([]byte, string) {
	hashGen := sha256.New()
	hashGen.Write([]byte(url))
	hash := hex.EncodeToString(hashGen.Sum(nil))

	return d.Get(hash)
}

func (d *DataStore) CacheData(url string, data []byte, mime string) {
	hashGen := sha256.New()
	hashGen.Write([]byte(url))
	hash := hex.EncodeToString(hashGen.Sum(nil))

	d.Set(hash, data, mime)
}

func (d *DataStore) GetPath(hash string) string {
	return d.Conf.StoragePath + "/" + hash
}
