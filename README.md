![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/i5heu/simple-S3-cache)
[![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/i5heu/simple-s3-cache)](https://hub.docker.com/repository/docker/i5heu/simple-s3-cache)
[![Docker Pulls](https://img.shields.io/docker/pulls/i5heu/simple-s3-cache)](https://hub.docker.com/repository/docker/i5heu/simple-s3-cache)

# simple-S3-cache
Just hook it in front of your public S3 bucket and enjoy reduction in bandwidth costs to your bucket

## Status
The project is working in its current state. It is not feature complete, but it is usable.
The code needs some cleanup, refactoring and most importantly tests.

## !WARNING!
This cache is still in development and is not ready for production use. It is not yet tested for security vulnerabilities.

## Available Environment Variables
* `SS3C_S3_ENDPOINT` - The endpoint of your S3 bucket. Defaults to `https://localhost`
* `SS3C_CORS_DOMAIN` - The domain to allow CORS requests from. Defaults to `https://example.com`
* `SS3C_USE_MAX_RAM_GB` - The maximum amount of RAM to use for caching. Defaults to `2` GB
* `SS3C_USE_MAX_DISK_GB` - The maximum amount of disk space to use for caching. Defaults to `25` GB
* `SS3C_STORAGE_PATH` - The path to store cached files. Defaults to `/cache`
* `SS3C_INFLUXDB_URL` - The URL of your InfluxDB instance. Defaults to `""`
* `SS3C_INFLUXDB_TOKEN` - The token to use for authentication with InfluxDB. Defaults to `""`
* `SS3C_INFLUXDB_ORG` - The organization to use for authentication with InfluxDB. Defaults to `""`
* `SS3C_INFLUXDB_BUCKET` - The bucket to use for authentication with InfluxDB. Defaults to `""`

## Docker Compose Example
```yaml
version: "3.7"
services:
  ss3c:
    image: i5heu/simple-s3-cache:latest
    container_name: ss3c
    restart: always
    environment:
      - SS3C_S3_ENDPOINT=https://cdn.example.com
      - SS3C_CORS_DOMAIN=https://example.com
      - SS3C_USE_MAX_RAM_GB=5
      - SS3C_USE_MAX_DISK_GB=50
      - SS3C_INFLUXDB_URL=http://influxdb.example.com:8086
      - SS3C_INFLUXDB_TOKEN=super-cool-token-iahfpwihg9h32ithowh
      - SS3C_INFLUXDB_ORG=super-cool-org
      - SS3C_INFLUXDB_BUCKET=super-cool-bucket
    ports:
      - 8080:8084
    volumes:
      - ./cache:/cache
```

## Docker Hub
https://hub.docker.com/repository/docker/i5heu/simple-s3-cache