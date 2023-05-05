package pinecone

import "fmt"

// buildFetchVectorPathParams builds the fetch vector path parameters.
// Example: ids=1&ids=2&ids=3&namespace=foo
func buildFetchVectorPathParams(params FetchVectorsParams) string {
	var pathParams string
	for i, id := range params.IDs {
		if i == len(params.IDs)-1 {
			pathParams += fmt.Sprintf("ids=%s", id)
			break
		}
		pathParams += fmt.Sprintf("ids=%s&", id)
	}
	if params.Namespace != "" {
		pathParams += fmt.Sprintf("&namespace=%s", params.Namespace)
	}
	return pathParams
}
