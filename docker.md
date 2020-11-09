# OpenBooks Docker Image

## Usage

### Basic

`docker run -d -p 8080:80 openbooks`

### Persist eBook Files

`docker run -d -p 8080 -v ~/Downloads:/books openbooks -persist`

### Set custom IRC name (default random adjective-username combination)

`docker run -d -p 8080:80 openbooks -name my_irc_name`

## Arguments

```
-name string
    Use a custom name when connecting to irchighway
-persist
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
        command: -persist
        image: evan-buss/openbooks:latest

volumes:
    booksVolume:
```

