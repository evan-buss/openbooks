# [v5.0.0] - 2023-10-08

## Added
- The browser application now supports multiple servers. Click the connection icon in the bottom left to select a server.
- Results Grid Changes:
  - Added fuzzy search
  - Can show and hide columns
  - New filter UI

## Removed
- No longer supports CLI mode. This mode was a holdover from the original version of OpenBooks, which was a command line application.
- `--server` and `--tls` flags are removed. The browser application allows switching between servers and toggling TLS.
  - OpenBooks can now connect to  the Undernet #bookz channel, as this uses the same IRC operator software as IRC Highway.

# [v4.5.0] - 2023-01-08

## Added
-  Use `--useragent/-u` flag to optionally specify the [UserAgent](https://en.wikipedia.org/wiki/Client-to-client_protocol#VERSION) reported to the IRC server. Default remains `OpenBooks v4.5.0`.

## Breaking 
- `--name/-n` flag **must** be specified when starting the application. OpenBooks will no longer generate a random `noun_adjective` username.
- Only a single connection to the IRC server will be made. Opening a second browser tab will show an error message.


