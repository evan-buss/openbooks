# OpenBooks Docker Image

> See [Github](https://github.com/evan-buss/openbooks) for more information.

## Usage

### Basic

`docker run -d -p 8080:80 evanbuss/openbooks`

### Persist eBook Files

`docker run -d -p 8080 -v ~/Downloads:/books evanbuss/openbooks --persist`

### Set custom IRC name (default is a random 'adjective-noun' combination)

`docker run -d -p 8080:80 evanbuss/openbooks --name my_irc_name`

### Host at a sub path behind a reverse proxy

`docker run -d -p 8080:80 -e BASE_PATH=/openbooks/ evanbuss/openbooks`

## Arguments

```
--name string
    Use a custom name when connecting to irchighway
--persist
    Keep book files in the download dir. Default is to delete after sending.
```

## Docker Compose

```docker
version: '3.3'
services:
    openbooks:
        ports:
            - '8080:80'
        volumes:
            - 'booksVolume:/books'
        restart: unless-stopped
        container_name: OpenBooks
        command: --persist
        environment:
          - BASE_PATH=/openbooks/
        image: evan-buss/openbooks:latest

volumes:
    booksVolume:
```
