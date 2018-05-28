package mapquest

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	Host      = "open.mapquestapi.com"
	UserAgent = "MapQuest Open Data API Google Go Client v0.1"
)

func apiURL(path ...string) *url.URL {
	u := new(url.URL)
	u.Scheme = "https"
	u.Host = Host
	u.Path = "/" + strings.Join(path, "/")

	return u
}

// Client is the entry point to all services of the MapQuest Open Data API.
// See https://developer.mapquest.com/documentation/open/ for details about
// what you can do with the MapQuest API.
type Client struct {
	httpClient *http.Client
	key        string
	log        *log.Logger
}

// NewClient creates a new client for accessing the MapQuest API. You need
// to specify your AppKey here.
func NewClient(key string) *Client {
	return &Client{
		key:        key,
		httpClient: http.DefaultClient,
	}
}

// SetHTTPClient allows the caller to specify a special http.Client for
// invoking the MapQuest API. If you do not specify a http.Client, the
// http.DefaultClient from net/http is used.
func (c *Client) SetHTTPClient(client *http.Client) {
	if client == nil {
		client = http.DefaultClient
	}
	c.httpClient = client
}

// HTTPClient returns the registered http.Client. Notice that nil can
// be returned here.
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

// SetLogger sets the logger to use when e.g. debugging requests.
// Set to nil to disable logging (the default).
func (c *Client) SetLogger(logger *log.Logger) {
	c.log = logger
}

// StaticMap gives access to the MapQuest static map API
// described here: https://developer.mapquest.com/documentation/open/static-map-api/v5/
func (c *Client) StaticMap() *StaticMapAPI {
	return &StaticMapAPI{c: c}
}
