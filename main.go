package main

import (
	"fmt"

	"github.com/i5heu/simple-S3-cache/config"
)

// proxy files from s3
// CORS
// cache on filesystem
// cache in memory

func main() {
	conf := config.GetValues()
	fmt.Println(conf)
}
