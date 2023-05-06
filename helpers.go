package pinecone

import (
	"net/url"
)

// buildFetchVectorPathParams builds the fetch vector path parameters.
// Example: ids=1&ids=2&ids=3&namespace=foo
func buildFetchVectorPathParams(params FetchVectorsParams) string {
	var pathParams url.Values
	for _, id := range params.IDs {
		pathParams.Add("ids", id)
	}
	if params.Namespace != "" {
		pathParams.Add("namespace", params.Namespace)
	}

	return pathParams.Encode()
}
