### Binary

1. Download the latest release for your platform from the [releases page](https://github.com/evan-buss/openbooks/releases).
2. Run the binary
   - Linux users may have to run `chmod +x [binary name]` to make it executable
3. `./openbooks --help`
   - This will display all possible configuration values and introduce the two modes; CLI or Server.

### Docker

- Basic config
  - `docker run -p 8080:80 evanbuss/openbooks`
- Config to persist all eBook files to disk
  - `docker run -p 8080:80 -v /home/evan/Downloads/openbooks:/books evanbuss/openbooks --persist`

### Setting the Base Path

OpenBooks server doesn't have to be hosted at the root of your webserver. The basepath value allows you to host it behind a reverse proxy. The base path value must have opening and closing forward slashes (default "/").

- Docker
  - `docker run -p 8080:80 -e BASE_PATH=/openbooks/ evanbuss/openbooks`
- Binary
  - `./openbooks server --basepath /openbooks/`

## Usage

For a complete list of features use the `--help` flags on all subcommands.
For example `openbooks cli --help or openbooks cli download --help`. There are
two modes; Server or CLI. In CLI mode you interact and download books through
a terminal interface. In server mode the application runs as a web application
that you can visit in your browser.

Double clicking the executable will open the UI in your browser. In the future it may use [webviews](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) to provide a "native-like" desktop application.
