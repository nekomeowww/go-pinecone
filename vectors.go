package pinecone

import (
	"bytes"
	"context"
	"fmt"
)

type Vector struct {
	Id       string         `json:"id"`
	Values   []float32      `json:"values"`
	Metadata map[string]any `json:"metadata"`
}

type UpsertVectorsParams struct {
	Vectors   []*Vector `json:"vectors"`
	Namespace string    `json:"namespace"`
}

type UpsertVectorsResponse struct {
	UpsertedCount int `json:"upserted_count"`
}

func (ic *IndexClient) UpsertVectors(ctx context.Context, params UpsertVectorsParams) (*UpsertVectorsResponse, error) {
	if len(params.Vectors) == 0 {
		return nil, fmt.Errorf("%w: vectors are required", ErrInvalidParams)
	}
	if params.Namespace == "" {
		return nil, fmt.Errorf("%w: namespace is required", ErrInvalidParams)
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

type UpdateVectorParams struct {
	Values       []float32      `json:"values"`
	SparseValues *SparseValues  `json:"sparse_values"`
	SetMetadata  map[string]any `json:"set_metadata"`
	Id           string         `json:"id"`
	Namespace    string         `json:"namespace"`
}

type SparseValues struct {
	Indices []int32   `json:"indices"`
	Values  []float32 `json:"values"`
}

func (ic *IndexClient) UpdateVector(ctx context.Context, params UpdateVectorParams) error {
	if params.Id == "" {
		return fmt.Errorf("%w: vector id is required", ErrInvalidParams)
	}
	if params.Namespace == "" {
		return fmt.Errorf("%w: namespace is required", ErrInvalidParams)
	}
	if len(params.Values) == 0 {
		return fmt.Errorf("%w: values are required", ErrInvalidParams)
	}
	if params.SparseValues == nil && len(params.SparseValues.Indices) == 0 && len(params.SparseValues.Values) == 0 {
		return fmt.Errorf("%w: sparse values object is required", ErrInvalidParams)
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

type QueryParams struct {
	Filter          map[string]any `json:"filter"`
	IncludeValues   bool           `json:"include_values"`
	IncludeMetadata bool           `json:"include_metadata"`
	Vector          []float32      `json:"vector"`
	SparseVector    *SparseVector  `json:"sparse_vector"`
	Namespace       string         `json:"namespace"`
	TopK            int            `json:"top_k"`
	Id              string         `json:"id"`
}

type SparseVector struct {
	Indices []int32   `json:"indices"`
	Values  []float32 `json:"values"`
}

type MatchingVector struct {
	Id           string         `json:"id"`
	Values       []float32      `json:"values"`
	Score        float32        `json:"score"`
	SparseValues *SparseValues  `json:"sparse_values"`
	Metadata     map[string]any `json:"metadata"`
}

type QueryResponse struct {
	Matches   []*MatchingVector `json:"matches"`
	Namespace string            `json:"namespace"`
}

func (ic *IndexClient) Query(ctx context.Context, params QueryParams) (*QueryResponse, error) {
	if params.Namespace == "" {
		return nil, fmt.Errorf("%w: namespace is required", ErrInvalidParams)
	}
	if len(params.Vector) == 0 && (params.SparseVector == nil || len(params.SparseVector.Indices) == 0 || len(params.SparseVector.Values) == 0) {
		return nil, fmt.Errorf("%w: vector or sparse vector is required", ErrInvalidParams)
	}
	if params.TopK == 0 {
		return nil, fmt.Errorf("%w: top k is required", ErrInvalidParams)
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
