package pinecone

import (
	"bytes"
	"context"
	"fmt"
)

// DescribeIndexStatsParams represents the parameters for a describe index stats request.
type DescribeIndexStatsParams struct {
	Filter map[string]any `json:"filter"`
}

// DescribeIndexStatsResponse represents the response from a describe index stats request.
type DescribeIndexStatsResponse struct {
	Namespaces       map[string]*VectorCount `json:"namespaces"`
	Dimensions       int64                   `json:"dimensions"`
	IndexFullness    float32                 `json:"indexFullness"`
	TotalVectorCount int64                   `json:"totalVectorCount"`
}

// VectorCount represents the number of vectors in a namespace.
type VectorCount struct {
	VectorCount int64 `json:"vectorCount"`
}

func (ic *IndexClient) DescribeIndexStats(ctx context.Context, params DescribeIndexStatsParams) (*DescribeIndexStatsResponse, error) {
	var respBody DescribeIndexStatsResponse
	resp, err := ic.reqClient.
		R().
		SetContentType("application/json").
		SetBody(params).
		SetSuccessResult(&respBody).
		SetContext(ctx).
		Post("/describe_index_stats")

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		buffer := new(bytes.Buffer)
		_, err := buffer.ReadFrom(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %s, status code: %d", ErrRequestFailed, buffer.String(), resp.StatusCode)
	}

	return &respBody, nil
}

// QueryParams represents the parameters for a query request.
// See https://docs.pinecone.io/reference/query for more information.
type QueryParams struct {
	Filter          map[string]any `json:"filter"`
	IncludeValues   bool           `json:"includeValues"`
	IncludeMetadata bool           `json:"includeMetadata"`
	Vector          []float32      `json:"vector"`
	SparseVector    *SparseVector  `json:"sparseVector"`
	Namespace       string         `json:"namespace"`
	TopK            int64          `json:"topK"`
	ID              string         `json:"id"`
}

// SparseVector represents a sparse vector.
type SparseVector struct {
	Indices []int32   `json:"indices"`
	Values  []float32 `json:"values"`
}

// Vector represents a scored vector.
type Vector struct {
	ID           string         `json:"id"`
	Score        float32        `json:"score,omitempty"`
	Values       []float32      `json:"values"`
	SparseValues *SparseVector  `json:"sparseValues,omitempty"`
	Metadata     map[string]any `json:"metadata"`
}

// QueryResponse represents the response from a query request.
type QueryResponse struct {
	Matches   []*Vector `json:"matches"`
	Namespace string    `json:"namespace"`
}

// Query performs a query request.
func (ic *IndexClient) Query(ctx context.Context, params QueryParams) (*QueryResponse, error) {
	if err := validateQueryParams(params); err != nil {
		return nil, err
	}

	var respBody QueryResponse
	resp, err := ic.reqClient.
		R().
		SetContentType("application/json").
		SetBody(params).
		SetSuccessResult(&respBody).
		SetContext(ctx).
		Post("/query")

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		buffer := new(bytes.Buffer)
		_, err := buffer.ReadFrom(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %s, status code: %d", ErrRequestFailed, buffer.String(), resp.StatusCode)
	}

	return &respBody, nil
}

// DeleteVectorsParams represents the parameters for a delete vectors request.
// See https://docs.pinecone.io/reference/delete_post for more information.
type DeleteVectorsParams struct {
	IDs       []string       `json:"ids"`
	Namespace string         `json:"namespace"`
	DeleteAll bool           `json:"deleteAll"`
	Filter    map[string]any `json:"filter"`
}

// DeleteVectors performs a delete vectors request.
// returns an error if the request fails and nil otherwise.
func (ic *IndexClient) DeleteVectors(ctx context.Context, params DeleteVectorsParams) error {
	if err := validateDeleteVectorsParams(params); err != nil {
		return err
	}

	resp, err := ic.reqClient.
		R().
		SetContentType("application/json").
		SetBody(params).
		SetContext(ctx).
		Post("/vectors/delete")

	if err != nil {
		return err
	}

	if !resp.IsSuccessState() {
		buffer := new(bytes.Buffer)
		_, err := buffer.ReadFrom(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%w: %s, status code: %d", ErrRequestFailed, buffer.String(), resp.StatusCode)
	}

	return nil
}

// FetchVectorsParams represents the parameters for a fetch vectors request.
type FetchVectorsParams struct {
	IDs       []string `json:"ids"`
	Namespace string   `json:"namespace"`
}

// FetchVectorsResponse represents the response from a fetch vectors request.
type FetchVectorsResponse struct {
	Vectors   map[string]*Vector `json:"vectors"`
	Namespace string             `json:"namespace"`
}

// FetchVectors performs a fetch vectors request.
// See https://docs.pinecone.io/reference/fetch for more information.
func (ic *IndexClient) FetchVectors(ctx context.Context, params FetchVectorsParams) (*FetchVectorsResponse, error) {
	if err := validateFetchVectorsParams(params); err != nil {
		return nil, err
	}

	pathParams := buildFetchVectorPathParams(params)
	var respBody FetchVectorsResponse

	resp, err := ic.reqClient.
		R().
		SetContentType("application/json").
		SetSuccessResult(&respBody).
		SetContext(ctx).
		Get("/vectors/fetch?" + pathParams)

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		buffer := new(bytes.Buffer)
		_, err := buffer.ReadFrom(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %s, status code: %d", ErrRequestFailed, buffer.String(), resp.StatusCode)
	}

	return &respBody, nil
}

// UpdateVectorParams represents the parameters for an update vector request.
// This has the same fields as the Vector struct plus setMetadata.
type UpdateVectorParams struct {
	Values       []float32      `json:"values"`
	SparseValues *SparseVector  `json:"sparseValues"`
	SetMetadata  map[string]any `json:"setMetadata"`
	ID           string         `json:"id"`
	Namespace    string         `json:"namespace"`
}

// UpdateVector performs an update vector request.
// See https://docs.pinecone.io/reference/update for more information.
func (ic *IndexClient) UpdateVector(ctx context.Context, params UpdateVectorParams) error {
	if err := validateUpdateVectorParams(params); err != nil {
		return err
	}

	resp, err := ic.reqClient.
		R().
		SetContentType("application/json").
		SetBody(params).
		SetContext(ctx).
		Post("/vectors/update")

	if err != nil {
		return err
	}

	if !resp.IsSuccessState() {
		buffer := new(bytes.Buffer)
		_, err := buffer.ReadFrom(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%w: %s, status code: %d", ErrRequestFailed, buffer.String(), resp.StatusCode)
	}

	return nil
}

// UpsertVectorsParams represents the parameters for an upsert vectors request.
// See https://docs.pinecone.io/reference/upsert for more information.
type UpsertVectorsParams struct {
	Vectors   []*Vector `json:"vectors"`
	Namespace string    `json:"namespace"`
}

// UpsertVectorsResponse represents the response from an upsert vectors request.
type UpsertVectorsResponse struct {
	UpsertedCount int `json:"upsertedCount"`
}

// UpsertVectors performs an upsert vectors request.
// See https://docs.pinecone.io/reference/upsert for more information.
func (ic *IndexClient) UpsertVectors(ctx context.Context, params UpsertVectorsParams) (*UpsertVectorsResponse, error) {
	if err := validateUpsertVectorsParams(params); err != nil {
		return nil, err
	}

	var respBody UpsertVectorsResponse
	resp, err := ic.reqClient.
		R().
		SetContentType("application/json").
		SetBody(params).
		SetSuccessResult(&respBody).
		SetContext(ctx).
		Post("/vectors/upsert")

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		buffer := new(bytes.Buffer)
		_, err := buffer.ReadFrom(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %s, status code: %d", ErrRequestFailed, buffer.String(), resp.StatusCode)
	}

	return &respBody, nil
}
