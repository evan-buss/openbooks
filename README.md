# openbooks
Openbooks allows you to download ebooks from irc.irchighway.net quickly and easily.

![home](https://raw.githubusercontent.com/evan-buss/openbooks/master/.github/home.png)
![search results](https://raw.githubusercontent.com/evan-buss/openbooks/master/.github/search.png)
## Getting Started
1. Download the latest release for your platform from the [releases page](https://github.com/evan-buss/openbooks/releases).
2. Run the binary
   - Linux users may have to run `chmod +x [binary name]` to make it executable

## Usage
There are two different usage modes. Either through the command line interface or the web interface

### Web UI (Recommended)
The web UI simplifies the process of searching and downloading books. It automatically downloads the results file, unzips
it and then parses the data so it can be searched and easily read. You are presented with a table of search results.
You can further refine your search by clicking the button at the top of each column. When you find what you are looking
for, click download and the book will be sent to you.

- Launch the application normally
  - `OpenBooks`
- Your web browser will be opened to the OpenBooks interface.
  - You must wait 30 seconds before your first search as there is a mandatory delay on the IRC server
  - On the left side you will see your past searches as well as a list of online download servers
- Enter your search query.
  - The book is sent to the backend server and placed in the queue
  - When the results are recieved they are then parsed and put into a      table for easy filtering.
  - The recommended button refines the search to servers that I 
    personally recommend and use whenever possible as they give 
    consistent results.
- When you find what you are looking for simply click the download
  button. 
  - It is recommended that you ensure the server is online before 
    clicking download.
- At any time you can select a past search from the sidebar to take 
  another look at results you may have missed.

### CLI 
- Launch the application with the -cli flag
  - `OpenBooks -cli`
- Search eBooks by entering `s` for search. Then enter your query.
  - A text file of search results is downloaded to your `~/Downloads` folder
  - Example Results
  ```
  !LawdyServer Kurt Vonnegut - Slaughterhouse Five.pdf  ::INFO:: 220.0KB
  !Oatmeal Vonnegut, Kurt - Slaughterhouse-Five(v1.1)[htm].rar  ::INFO:: 257.0KB 
  !Oatmeal 023 - Vonnegut, Kurt Jr. - Slaughterhouse-five (v1.1) [html, jpg].rar  ::INFO:: 257.0KB 
  !Oatmeal Kurt Vonnegut - Slaughterhouse-Five (v4.0) (html).rar  ::INFO:: 245.9KB 
  !Oatmeal Kurt Vonnegut - Slaughterhouse-Five (v5.0) (epub).rar  ::INFO:: 1.9MB 
  !Pondering42 023 - Vonnegut, Kurt Jr. - Slaughterhouse-five (v1.1) [html, jpg].rar ::INFO:: 256.99KB
  ``` 
- Download eBooks by entering `d` for download. Then paste in your book string
  - Enter the download string of the book you want to receive (don't include the ::INFO:: )
    - Ex) `!Oatmeal Kurt Vonnegut - Slaughterhouse-Five (v5.0) (epub).rar`
  - The eBook file will be downloaded to your `~/Downloads` directory
  - If the book is contained within a compressed archive such as .zip or .rar, it is extracted automatically

### Optional Launch Flags
```
 -cli 
    Launch OpenBooks in the terminal instead of the web UI
 -log 
    Save IRC logs to "irc_log.txt"
 -name string 
    Use a name that differs from your account name. (default "[os username]")
```

## Development

### Install the dependencies
  - `go get`
  - `go get github.com/rakyll/statik`
    - This is installed as a binary to your `$GOBIN`
  - `cd server/app && npm install`

### Build the React SPA and compile the complete binary
  - `cd server/app && npm run build`
  - Go back to root directory: `./statik.sh && go build`

### Build the go binary (if you haven't changed the frontend)
  - `go build`

### Build the complete binary for all platforms
  - `./build.sh`

## Why / How
- I wrote this as an easier way to search and download books from irchighway.net. It handles all the extraction and data processing for you. You just have to click the book you want. Hopefully you find it much easier than the IRC interface.
- It was also interesting to learn how the [IRC](https://en.wikipedia.org/wiki/Internet_Relay_Chat) and [DCC](https://en.wikipedia.org/wiki/Direct_Client-to-Client) protocols work and write custom implementations.

## Technology

- Backend
  - Golang
  - Packr2 (bundle static assets in the go binary)
  - Archiver (extract files from archive)
  - gorilla/websocket (communication between backend and UI)
- Frontend
  - React.js
  - Ant Design Components
