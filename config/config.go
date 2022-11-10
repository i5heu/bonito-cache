package config

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
)

// https://masto-cdn.s3.eu-central-003.backblazeb2.com/accounts/avatars/109/292/184/771/206/331/original/79ff3fa20a3d602d.gif

type Config struct {
	S3Endpoint   string `example:"https://localhost" env:"SS3C_S3_ENDPOINT"`   //has trailing slash
	CORSDomain   string `example:"https://example.com" env:"SS3C_CORS_DOMAIN"` //has trailing slash
	UseMaxRamGb  int    `example:"2" env:"SS3C_USE_MAX_RAM_GB"`
	UseMaxDiskGb int    `example:"25" env:"SS3C_USE_MAX_DISK_GB"`
}

func GetValues() Config {
	c := Config{}
	fields := reflect.VisibleFields(reflect.TypeOf(struct{ Config }{}))

	for _, field := range fields {
		switch field.Type {
		case reflect.TypeOf(string("")):
			urlCleaned := checkAndCleanURL(getEnvValueString(field.Tag.Get("env")))

			if urlCleaned == "" {
				urlCleaned = field.Tag.Get("example")
			}

			reflect.ValueOf(&c).Elem().FieldByName(field.Name).SetString(urlCleaned)

		case reflect.TypeOf(int(0)):
			intValue := getEnvValueInt(field.Tag.Get("env"))
			if intValue == 0 {
				var err error
				intValue, err = strconv.Atoi(field.Tag.Get("example"))
				if err != nil {
					panic(err)
				}
			}
			reflect.ValueOf(&c).Elem().FieldByName(field.Name).SetInt(int64(intValue))
		}
	}

	return c
}

func checkAndCleanURL(urlDirty string) string {
	urlCleaned := urlDirty

	if urlCleaned == "" {
		return urlCleaned
	}

	// check if url is valid
	_, err := url.ParseRequestURI(urlCleaned)
	if err != nil {
		panic(err)
	}

	// check if url has trailing slash
	if urlCleaned[len(urlCleaned)-1:] != "/" {
		urlCleaned = urlCleaned + "/"
	}

	// check if url has protocol
	if urlCleaned[0:4] != "http" {
		urlCleaned = "https://" + urlCleaned
	}

	return urlCleaned
}

func getEnvValueString(env string) string {
	fmt.Println("---!>", os.Getenv(env), env)
	return os.Getenv(env)
}
func getEnvValueInt(env string) int {
	foo := os.Getenv(env)
	fmt.Println(foo, "--", env)
	if foo == "" {
		return 0
	}

	//  string to int
	bar, err := strconv.Atoi(foo)
	if err != nil {
		panic(err)
	}

	return bar
}

func GetCompleteURL(c Config, path string) string {
	return c.S3Endpoint + path
}
