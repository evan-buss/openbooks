# [v4.5.0] - 2023-01-08

## Added
-  Use `--useragent/-u` flag to optionally specify the [UserAgent](https://en.wikipedia.org/wiki/Client-to-client_protocol#VERSION) reported to the IRC server. Default remains `OpenBooks v4.5.0`.

## Breaking 
- `--name/-n` flag **must** be specified when starting the application. OpenBooks will no longer generate a random `noun_adjective` username.
- Only a single connection to the IRC server will be made. Opening a second browser tab will show an error message.


