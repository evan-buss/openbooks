OpenBooks provides a convenient user interface over [IRC Highway's](https://irchighway.net/) `#ebook` channel.
It streamlines the process of connecting, searching for, and downloading books.

!!! tip "OpenBooks does not host any content, think of it as a single-purpose IRC client."

There are 2 _modes of operation_; CLI or Server.

The majority of users will want to check out the [Server Mode](./configuration.md) where you search and download via a web interface in your browser.
This allows you to self-host OpenBooks without having to install it on every device.

If you'd prefer to use OpenBooks from your terminal, check out [CLI Mode](./configuration.md).

## Organized Downloads

You can opt-in to save downloaded books in a structured directory:

    books/Author Name/Book Title/FILE

This makes managing your library easier. Enable this with the `--organize-downloads` flag. By default, all downloads go into a single directory.

### Docker

`docker run -p 8080:80 evanbuss/openbooks`

: Basic configuration that exposes the web interface on [http://localhost:8080](http://localhost:8080) and saves all files to an anonymous volume.

`docker run -p 8080:80 -v ~/Downloads/openbooks:/books evanbuss/openbooks --persist`

: More advanced configuration that exposes the web interface on [http://localhost:8080](http://localhost:8080) and persists all eBook files to the mounted volume at `~/Downloads/openbooks`.

> For more information see the [docker guide](./setup/docker.md).

### Executable

1. Download the latest release for your platform from the [releases page](https://github.com/evan-buss/openbooks/releases).
2. Execute it from your terminal in Server (`./openbooks server`) or CLI (`./openbooks cli`) mode.

   - Linux users may have to run `chmod +x [binary name]` to make it executable

> For more information see the [executable guide](./setup/executable.md).
