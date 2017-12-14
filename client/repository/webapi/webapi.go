package webapi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/slok/ragnarok/api"
	apiutil "github.com/slok/ragnarok/api/util"
	"github.com/slok/ragnarok/apimachinery/serializer"
	"github.com/slok/ragnarok/apimachinery/watch"
	"github.com/slok/ragnarok/log"
)

// Client implements the access to store and retrieve resources from Ragnarok's rest HTTP API. Satisfies repository.Client interface.
type Client struct {
	httpCli      *http.Client
	serializer   serializer.Serializer
	resourceType api.TypeMeta
	url          *url.URL
	logger       log.Logger
}

// NewDefaultClient returns a new Ragnarok API client with some defaults.
func NewDefaultClient(apiURL string, httpCli *http.Client, resourceType api.TypeMeta, logger log.Logger) (*Client, error) {
	return NewClient(apiURL, httpCli, resourceType, serializer.DefaultSerializer, logger)
}

// NewClient returns a new Ragnarok API client.
func NewClient(apiURL string, httpCli *http.Client, resourceType api.TypeMeta, serializer serializer.Serializer, logger log.Logger) (*Client, error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		httpCli:      httpCli,
		serializer:   serializer,
		resourceType: resourceType,
		url:          u,
		logger:       logger,
	}, nil
}

func (c *Client) getAPIURL(id string) string {
	u := fmt.Sprintf("%s%s", c.url.String(), apiutil.GetTypeAPIPath(c.resourceType))
	if id != "" {
		u = fmt.Sprintf("%s/%s", u, id)
	}
	return u
}

func (c *Client) setContentType(r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
}

// Create will create an object using the HTTP rest API. Satisfies repository.Client interface.
func (c *Client) Create(obj api.Object) (api.Object, error) {
	// Serialize the object.
	var b bytes.Buffer
	if err := c.serializer.Encode(obj, &b); err != nil {
		return nil, err
	}

	// Create the request.
	u := c.getAPIURL("")
	req, err := http.NewRequest("POST", u, &b)
	if err != nil {
		return nil, err
	}
	c.setContentType(req)

	// Make the call.
	resp, err := c.httpCli.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("server return an %d status code: %s", resp.StatusCode, body)
	}

	return c.serializer.Decode(body)
}

// Update will update an object using the HTTP rest API. Satisfies repository.Client interface.
func (c *Client) Update(obj api.Object) (api.Object, error) {
	return nil, fmt.Errorf("not implemented")
}

// Delete will Delete an object using the HTTP rest API. Satisfies repository.Client interface.
func (c *Client) Delete(id string) error {
	return fmt.Errorf("not implemented")
}

// Get will get an object using the HTTP rest API. Satisfies repository.Client interface.
func (c *Client) Get(id string) (api.Object, error) {
	return nil, fmt.Errorf("not implemented")
}

// List will list objects using the HTTP rest API. Satisfies repository.Client interface.
func (c *Client) List(opts api.ListOptions) ([]api.Object, error) {
	return nil, fmt.Errorf("not implemented")
}

// Watch will watch an object type using the HTTP rest API. Satisfies repository.Client interface.
func (c *Client) Watch(opts api.ListOptions) (watch.Watcher, error) {
	return nil, fmt.Errorf("not implemented")
}
