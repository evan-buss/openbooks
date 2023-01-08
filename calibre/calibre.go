package calibre

import (
	"encoding/json"
	"github.com/icholy/digest"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync/atomic"
)

type Client struct {
	*http.Client
	config *Config
	jobId  uint32
}

func NewClient(config *Config) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	httpClient := http.DefaultClient
	if config.Username != "" && config.Password != "" {
		httpClient = &http.Client{
			Transport: &digest.Transport{
				Username: config.Username,
				Password: config.Password,
			},
		}
	}

	client := &Client{
		Client: httpClient,
		config: config,
	}

	return client, nil
}

func (client *Client) AddBook(bookPath string) (AddBookResponse, error) {

	duplicates := "n"
	if client.config.AddDuplicates {
		duplicates = "y"
	}

	// Increment Job ID (not really sure what this is used for)
	atomic.AddUint32(&client.jobId, 1)
	//Example Request /cdb/add-book/1/n/The Great Gatsby.epub
	uri, err := url.Parse(client.config.URL + path.Join("/cdb/add-book", strconv.Itoa(int(client.jobId)), duplicates, filepath.Base(bookPath)))
	if err != nil {
		return AddBookResponse{}, err
	}

	file, err := os.Open(bookPath)
	defer file.Close()
	if err != nil {
		return AddBookResponse{}, err
	}

	resp, err := client.Post(uri.String(), "application/epub+zip", file)
	if err != nil {
		return AddBookResponse{}, err
	}

	defer resp.Body.Close()

	var response AddBookResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return AddBookResponse{}, err
	}

	return response, nil
}
