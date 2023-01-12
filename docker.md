# OpenBooks Docker Image

> See [Github](https://github.com/evan-buss/openbooks) for more information.

## Usage

### Basic

`docker run -d -p 8080:80 evanbuss/openbooks --name my_irc_name`

### Persist eBook Files

`docker run -d -p 8080 -v ~/Downloads:/books evanbuss/openbooks --name my_irc_name --persist`

### Host at a sub path behind a reverse proxy

`docker run -d -p 8080:80 -e BASE_PATH=/openbooks/ evanbuss/openbooks --name my_irc_name`

## Arguments

```
--name string
    Required name when connecting to irchighway
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
        command: --name my_irc_name --persist
        environment:
          - BASE_PATH=/openbooks/
        image: evanbuss/openbooks:latest

volumes:
    booksVolume:
```
