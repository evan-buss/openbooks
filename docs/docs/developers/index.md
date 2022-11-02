# Getting Started

OpenBooks makes use of the [Task](https://taskfile.dev/) build system which provides a cross platform and simple alternative to GNU Make.

If making use of the [devcontainer](./dev-container.md), Task will automatically be set up for you.
Otherwise follow the install instructions on the Task website. The rest of the guide assumes you have Task installed.

`task dev:init`

: Installs NPM and Go dependencies.

`task dev:mock`

: Starts a mock IRC / DCC server that mimics basic requests / responses from IRC Highway. 95% of the time you should be connecting to the mock server to avoid making requests to the IRC Highway's server. See [IRC notes](../irc-notes.md).

## Server Mode Development

Run the following commands in separate terminals.

`task dev:client`

: Starts the React front-end and enabled hot reload via Vite.

`task dev:server`

: Compiles and runs OpenBooks in Server mode. Connects to the Mock IRC server.

## CLI Mode Development

`task dev:cli`

: Compiles and runs OpenBooks in CLI mode. Connects to the Mock IRC server.

<!-- ## Why / How

- I wrote this as an easier way to search and download books from irchighway.net. It handles all the extraction and data processing for you. You just have to click the book you want. Hopefully you find it much easier than the IRC interface.
- It was also interesting to learn how the [IRC](https://en.wikipedia.org/wiki/Internet_Relay_Chat) and [DCC](https://en.wikipedia.org/wiki/Direct_Client-to-Client) protocols work and write custom implementations.

## Technology

- Backend
  - Golang
  - Chi
  - gorilla/websocket
  - Archiver (extract files from various archive formats)
- Frontend
  - React.js
  - TypeScript
  - Redux / Redux Toolkit
  - Mantine UI / @emotion/react
  - Framer Motion -->
