[![Docker Pulls](https://img.shields.io/docker/pulls/evanbuss/openbooks.svg)](https://hub.docker.com/r/
evanbuss/openbooks/)

The OpenBooks docker image allows you to run [Server Mode](../modes/server.md). A multi-platform Docker container is published to Docker Hub for each release.

## Docker Compose

For advanced configuration, I recommend using Docker Compose to keep track of container setup.

```yaml title="docker-compose.yml"
version: "3.3"
services:
  openbooks:
    container_name: OpenBooks
    image: evanbuss/openbooks:latest
    restart: unless-stopped
    ports:
      - "8080:80"
    volumes:
      - "~/Downoads/openbooks:/books"
    command: --persist --name
    environment:
      - BASE_PATH=/openbooks/
```

## Configuration

See the [configuration docs](../configuration.md) for a complete list of Server mode configuration options. Pass the configuration flags into the `command` property.

Use the `environment` property to optionally set a custom base path for the server.

## Image Tags

`evanbuss/openbooks:latest`

: The majority of users will want this image and will always be up to date with the latest release. Note that auto-updating between version could break configuration.[^1]

`evanbuss/openbooks:X.X.X`

: Version specific tags. Each time a new release is cut, a new version tagged image is published.

`evanbuss/openbooks:edge`

: Built from the latest development build. This image is best if you want to test the latest changes but be warned that it could be unstable and not work at all.

## Image Platforms

- `linux/amd64`
- `linux/arm64`
- `linux/arm`

[^1]: I personally auto-update all of my docker containers and haven't experienced many issues. Tools like [Watchtower](https://containrrr.dev/watchtower/) can check for updates, pull images, and restart containers automatically.
