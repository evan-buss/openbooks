## [Development]

### Added

- The IRC connection is kept open between page reloads / websocket disconnects.
    - After the websocket is disconnected, the IRC connection will remain open until the websocket is reconnected or 3 minutes have passed.
      This will prevent inadvertent bans from the IRC channel due to frequent connects/disconnects.
- Grid Improvements
    - Column sorting
    - Fuzzy searching across title / author columns
    - Toggle column visibility

## [v4.5.0] - 2023-01-08

### Added

- Use `--useragent/-u` flag to optionally specify the [UserAgent](https://en.wikipedia.org/wiki/Client-to-client_protocol#VERSION) reported to the IRC server. Default remains `OpenBooks v4.5.0`.

### Breaking

- `--name/-n` flag **must** be specified when starting the application. OpenBooks will no longer generate a random `noun_adjective` username.
- Only a single connection to the IRC server will be made. Opening a second browser tab will show an error message.
