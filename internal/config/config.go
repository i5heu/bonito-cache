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
	S3Endpoint     string `kind:"url" example:"https://localhost" env:"BONITO_S3_ENDPOINT"`   //has trailing slash
	CORSDomain     string `kind:"url" example:"https://example.com" env:"BONITO_CORS_DOMAIN"` //has trailing slash
	UseMaxRamGB    int    `kind:"int" example:"2" env:"BONITO_USE_MAX_RAM_GB"`
	UseMaxDiskGb   int    `kind:"int" example:"25" env:"BONITO_USE_MAX_DISK_GB"`
	StoragePath    string `kind:"string" example:"/cache" env:"BONITO_STORAGE_PATH"`
	InfluxDbUrl    string `kind:"url" example:"" env:"BONITO_INFLUXDB_URL"`
	InfluxDbToken  string `kind:"string" example:"" env:"BONITO_INFLUXDB_TOKEN"`
	InfluxDbOrg    string `kind:"string" example:"" env:"BONITO_INFLUXDB_ORG"`
	InfluxDbBucket string `kind:"string" example:"" env:"BONITO_INFLUXDB_BUCKET"`
	ClusterActive  bool   `kind:"bool" example:"false" env:"BONITO_CLUSTER_ACTIVE"`
	ClusterSeed    string `kind:"url" example:"https://node1.example.com" env:"BONITO_CLUSTER_SEED"` // must be https!
	ClusterKey     string `kind:"string" example:"" env:"BONITO_CLUSTER_KEY"`
}

func GetValues() Config {
	c := Config{}
	fields := reflect.VisibleFields(reflect.TypeOf(struct{ Config }{}))

	for _, field := range fields {
		switch field.Tag.Get("kind") {
		case "url":
			urlCleaned := checkAndCleanURL(getEnvValueString(field.Tag.Get("env")))

			if urlCleaned == "" {
				urlCleaned = field.Tag.Get("example")
			}

			reflect.ValueOf(&c).Elem().FieldByName(field.Name).SetString(urlCleaned)
		case "string":
			value := getEnvValueString(field.Tag.Get("env"))
			if value == "" {
				value = field.Tag.Get("example")
			}

			reflect.ValueOf(&c).Elem().FieldByName(field.Name).SetString(value)
		case "bool":
			value := getEnvValueString(field.Tag.Get("env"))
			switch value {
			case "true":
				reflect.ValueOf(&c).Elem().FieldByName(field.Name).SetBool(true)
			case "":
			case "false":
				reflect.ValueOf(&c).Elem().FieldByName(field.Name).SetBool(false)
			default:
				fmt.Println("Error in env variable " + field.Name + " Invalid value for bool: " + value)
				reflect.ValueOf(&c).Elem().FieldByName(field.Name).SetBool(false)
			}

		case "int":
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

	// remove trailing slash if present
	if urlCleaned[len(urlCleaned)-1:] == "/" {
		urlCleaned = urlCleaned[:len(urlCleaned)-1]
	}

	// check if url has protocol
	if urlCleaned[0:4] != "http" {
		urlCleaned = "https://" + urlCleaned
	}

	return urlCleaned
}

func getEnvValueString(env string) string {
	return os.Getenv(env)
}
func getEnvValueInt(env string) int {
	foo := os.Getenv(env)
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
