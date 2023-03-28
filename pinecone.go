package pinecone

import (
	"fmt"

	"github.com/imroc/req/v3"
)

// Client is the main entry point for the Pinecone API.
type Client struct {
	options   *options
	reqClient *req.Client
}

// New creates a new Pinecone client.
func New(callOpts ...CallOptions) (*Client, error) {
	opts := applyCallOptions(callOpts)
	reqClient := req.
		C().
		SetBaseURL(fmt.Sprintf("https://controller.%s.pinecone.io", opts.environment)).
		SetCommonHeader("Api-Key", opts.apiKey)
	return &Client{
		options:   opts,
		reqClient: reqClient,
	}, nil
}

// Debug enables debug logging and http dump for the client.
func (c *Client) Debug() *Client {
	c.reqClient.DebugLog = true
	c.reqClient = c.reqClient.EnableDumpAll()
	return c
}
