package ramCache

import (
	"fmt"
	"runtime"
	"sort"
	"time"

	"github.com/inhies/go-bytesize"
)

type File struct {
	Hash string
	Size uint // in bytes
	hits uint
}

func (d *DataStore) RamFileManager() {

	go func() {
		for {
			time.Sleep(5 * time.Second)
			fmt.Println("garbage collector running")
			d.Ch <- File{Hash: "GC", hits: 42}
		}
	}()

	hitMap := make(map[string]File)

	for file := range d.Ch {
		if file.Hash == "GC" {
			d.garbageCollector(hitMap)
		} else {
			if hashMapFile, ok := hitMap[file.Hash]; ok {
				hashMapFile.hits++
			} else {
				hitMap[file.Hash] = file
			}
		}
	}
}

func (d *DataStore) garbageCollector(hitMap map[string]File) {
	cacheSize := calculateCacheSize(hitMap)
	fmt.Println("cache size:", bytesize.ByteSize(cacheSize).Format("%.5f", "GB", false))
	timeStart := time.Now()

	if cacheSize < uint(d.Conf.UseMaxRamGB*int(bytesize.GB)) {
		return
	}

	fileSizeRemoved := uint(0)
	filesRemoved := 0
	for _, file := range sortHitMapByHits(hitMap) {
		if cacheSize < uint(d.Conf.UseMaxRamGB*int(bytesize.GB)) {
			fmt.Println("Removed:", filesRemoved, "files with a total size of:", bytesize.ByteSize(fileSizeRemoved).Format("%.5f", "GB", false))
			d.Log.LogCache(timeStart, "ramGC", cacheSize, uint(d.Conf.UseMaxRamGB*int(bytesize.GB)))
			runtime.GC()
			return
		}

		d.Data.Delete(file.Hash)
		delete(hitMap, file.Hash)
		cacheSize -= file.Size
		fileSizeRemoved += file.Size
		filesRemoved++
	}
}

func calculateCacheSize(hitMap map[string]File) uint {
	var size uint
	for _, file := range hitMap {
		size += file.Size
	}
	return size
}

func sortHitMapByHits(hitMap map[string]File) []File {
	var files []File
	for _, file := range hitMap {
		files = append(files, file)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].hits < files[j].hits
	})

	return files
}
