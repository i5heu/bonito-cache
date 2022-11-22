![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/i5heu/bonito-cache)
[![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/i5heu/simple-S3-cache)](https://hub.docker.com/repository/docker/i5heu/simple-S3-cache)
[![Docker Pulls](https://img.shields.io/docker/pulls/i5heu/simple-S3-cache)](https://hub.docker.com/repository/docker/i5heu/simple-S3-cache)
[![wakatime](https://wakatime.com/badge/github/i5heu/bonito-cache.svg)](https://wakatime.com/badge/github/i5heu/bonito-cache)

![Logo of a humanoid fish holding notes and a note box](./media/logo_small.png)

# bonito-cache
Just hook it in front of your public S3 bucket and enjoy reduction in bandwidth costs to your bucket

## Status
The project is working in its current state. It is not feature complete, but it is usable.
The code needs some cleanup, refactoring and most importantly tests.

## !WARNING!
This cache is still in development and is not ready for production use. It is not yet tested for security vulnerabilities.

## Available Environment Variables
* `BONITO_S3_ENDPOINT` - The endpoint of your S3 bucket. Defaults to `https://localhost`
* `BONITO_CORS_DOMAIN` - The domain to allow CORS requests from. Defaults to `https://example.com`
* `BONITO_USE_MAX_RAM_GB` - The maximum amount of RAM to use for caching. Defaults to `2` GB
* `BONITO_USE_MAX_DISK_GB` - The maximum amount of disk space to use for caching. Defaults to `25` GB
* `BONITO_STORAGE_PATH` - The path to store cached files. Defaults to `/cache`
* `BONITO_INFLUXDB_URL` - The URL of your InfluxDB instance. Defaults to `""`
* `BONITO_INFLUXDB_TOKEN` - The token to use for authentication with InfluxDB. Defaults to `""`
* `BONITO_INFLUXDB_ORG` - The organization to use for authentication with InfluxDB. Defaults to `""`
* `BONITO_INFLUXDB_BUCKET` - The bucket to use for authentication with InfluxDB. Defaults to `""`

## Docker Compose Example
```yaml
version: "3.7"
services:
  ss3c:
    image: i5heu/bonito-cache:latest
    container_name: ss3c
    restart: always
    environment:
      - BONITO_S3_ENDPOINT=https://cdn.example.com
      - BONITO_CORS_DOMAIN=https://example.com
      - BONITO_USE_MAX_RAM_GB=5
      - BONITO_USE_MAX_DISK_GB=50
      - BONITO_INFLUXDB_URL=http://influxdb.example.com:8086
      - BONITO_INFLUXDB_TOKEN=super-cool-token-iahfpwihg9h32ithowh
      - BONITO_INFLUXDB_ORG=super-cool-org
      - BONITO_INFLUXDB_BUCKET=super-cool-bucket
    ports:
      - 8080:8084
    volumes:
      - ./cache:/cache
```

## Docker Hub
https://hub.docker.com/repository/docker/i5heu/bonito-cache

## TODO for 1.0.0
- Add a CLI and API to delete a cached file
- refactor code
- build better docs
- rename project

## Future features 

Caching in cluster mode.  
If you build a cluster of bonito-cache instances, they will try to request a file from each other first before requesting it from the S3 bucket.  

## License
bonito-cache © 2022 Mia Heidenstedt and contributors   
SPDX-License-Identifier: AGPL-3.0  
