package pinecone

import (
	"fmt"
	"github.com/imroc/req/v3"
)

// IndexClient client for vector operations
type IndexClient struct {
	reqClient *req.Client
}

func NewIndexClient(opts ...CallOptions) *IndexClient {
	appliedOptions := applyCallOptions(opts)
	reqClient := req.
		C().
		SetBaseURL(fmt.Sprintf("https://%s-%s.svc.%s.pinecone.io", appliedOptions.indexName, appliedOptions.projectName, appliedOptions.environment)).
		SetCommonHeader("Api-Key", appliedOptions.apiKey)
	return &IndexClient{
		reqClient: reqClient,
	}
}

func (c *IndexClient) Debug() *IndexClient {
	c.reqClient.DebugLog = true
	c.reqClient = c.reqClient.EnableDumpAll()
	return c
}
