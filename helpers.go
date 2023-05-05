package pinecone

func buildFetchVectorPathParams(params FetchVectorsParams) map[string]string {
	pathParams := map[string]string{}
	for _, id := range params.Ids {
		pathParams["ids"] = id
	}
	if params.Namespace != "" {
		pathParams["namespace"] = params.Namespace
	}
	return pathParams
}
