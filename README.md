# openbooks

[![Docker Pulls](https://img.shields.io/docker/pulls/evanbuss/openbooks.svg)](https://hub.docker.com/r/evanbuss/openbooks/)

Openbooks allows you to download ebooks from irc.irchighway.net quickly and easily.

![home](https://raw.githubusercontent.com/evan-buss/openbooks/master/.github/home.png)
![search results](https://raw.githubusercontent.com/evan-buss/openbooks/master/.github/search.png)

## Getting Started

### Binary

1. Download the latest release for your platform from the [releases page](https://github.com/evan-buss/openbooks/releases).
2. Run the binary
   - Linux users may have to run `chmod +x [binary name]` to make it executable
3. `./openbooks --help`
   - This will display all possible configuration values and introduce the two modes; CLI or Server.

### Docker

- Basic config
  - `docker run -p 8080:80 evanbuss/openbooks`
- Config to perist all eBook files to disk
  - `docker run -p 8080:80 -v /home/evan/Downloads/books:/books evanbuss/openbooks --persist`

## Development

### Install the dependencies

- `go get`
- `cd server/app && npm install`
- `cd ../..`
- `go run main.go`

### Build the React SPA and compile binaries for multiple platforms.

- Run `./build.sh`
- This will install npm packages, build the React app, and compile the executa

### Build the go binary (if you haven't changed the frontend)

- `go build`

## Why / How

- I wrote this as an easier way to search and download books from irchighway.net. It handles all the extraction and data processing for you. You just have to click the book you want. Hopefully you find it much easier than the IRC interface.
- It was also interesting to learn how the [IRC](https://en.wikipedia.org/wiki/Internet_Relay_Chat) and [DCC](https://en.wikipedia.org/wiki/Direct_Client-to-Client) protocols work and write custom implementations.

## Technology

- Backend
  - Golang
  - Archiver (extract files from various archive formats)
  - gorilla/websocket
- Frontend
  - React.js
  - TypeScript
  - Redux / Redux Toolkit
  - Styled Components
  - Evergreen UI
