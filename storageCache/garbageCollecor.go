package storageCache

import (
	"fmt"
	"io/ioutil"
	"sort"
	"time"

	"github.com/inhies/go-bytesize"
)

type File struct {
	Hash       string
	Size       uint // in bytes
	LastAccess time.Time
}

func (d *DataStore) StorageFileManager() {

	for {
		time.Sleep(5 * time.Second)
		d.garbageCollector()
	}
}

func (d *DataStore) garbageCollector() {
	cachePath := d.Conf.StoragePath
	fileStats := []File{}
	files, _ := ioutil.ReadDir(cachePath)

	for _, f := range files {
		f := File{Hash: f.Name(), Size: uint(f.Size()), LastAccess: f.ModTime()}
		fileStats = append(fileStats, f)
	}

	cacheSize := calculateCacheSize(fileStats)
	fmt.Println("cache size Storage:", bytesize.ByteSize(cacheSize).Format("%.5f", "GB", false))

	if cacheSize < uint(d.Conf.UseMaxDiskGb*int(bytesize.GB)) {
		return
	}

	// sort files by last access
	sort.Slice(fileStats, func(i, j int) bool {
		return fileStats[i].LastAccess.Before(fileStats[j].LastAccess)
	})

	fileSizeRemoved := uint(0)
	filesRemoved := 0
	for _, file := range fileStats {
		if cacheSize < uint(d.Conf.UseMaxDiskGb*int(bytesize.GB)) {
			fmt.Println("Removed from Storage:", filesRemoved, "files with a total size of:", bytesize.ByteSize(fileSizeRemoved).Format("%.5f", "GB", false))
			d.Log.LogCache(time.Now(), "storageGC", cacheSize, uint(d.Conf.UseMaxDiskGb*int(bytesize.GB)))
			return
		}
		d.delete(file.Hash)
		cacheSize -= file.Size
		fileSizeRemoved += file.Size
		filesRemoved++
	}
}

func calculateCacheSize(fileStats []File) uint {
	cacheSize := uint(0)
	for _, f := range fileStats {
		cacheSize += f.Size
	}
	return cacheSize
}
