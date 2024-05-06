## [Development]

### Added

- IRC Server Selection
    - Click the connection icon in the bottom left to select a server.
    - There are currently only 2 known IRC servers that support the exact same search / download syntax.
- The IRC connection is kept open between page reloads / websocket disconnects.
    - After the websocket is disconnected, the IRC connection will remain open until the websocket is reconnected or 3 minutes have passed.
      This will prevent inadvertent bans from the IRC channel due to frequent connects/disconnects.
- Grid Improvements
    - Column sorting
    - Fuzzy searching across title / author columns
    - Toggle column visibility
    - New filter UI

## Removed
- No longer supports CLI mode. This mode was a holdover from the original version of OpenBooks, which was a command line application.
- `--server` and `--tls` flags are removed. The browser application allows switching between servers and toggling TLS.
  - OpenBooks can now connect to  the Undernet #bookz channel, as this uses the same IRC operator software as IRC Highway.

# [v4.5.0] - 2023-01-08

### Added

- Use `--useragent/-u` flag to optionally specify the [UserAgent](https://en.wikipedia.org/wiki/Client-to-client_protocol#VERSION) reported to the IRC server. Default remains `OpenBooks v4.5.0`.

### Breaking

- `--name/-n` flag **must** be specified when starting the application. OpenBooks will no longer generate a random `noun_adjective` username.
- Only a single connection to the IRC server will be made. Opening a second browser tab will show an error message.
