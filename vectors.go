package pinecone

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

type Vector struct {
	Id       string         `json:"id"`
	Values   []float32      `json:"values"`
	Metadata map[string]any `json:"metadata"`
}

type UpsertVectorsParams struct {
	Vectors   []*MatchingVector `json:"vectors"`
	Namespace string            `json:"namespace"`
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
	SparseValues *SparseValues  `json:"sparseValues"`
	SetMetadata  map[string]any `json:"setMetadata"`
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
	IncludeValues   bool           `json:"includeValues"`
	IncludeMetadata bool           `json:"includeMetadata"`
	Vector          []float32      `json:"vector"`
	SparseVector    *SparseVector  `json:"sparseVector"`
	Namespace       string         `json:"namespace"`
	TopK            int            `json:"topK"`
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
	SparseValues *SparseValues  `json:"sparseValues"`
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

type FetchVectorsParams struct {
	Ids       []string `json:"ids"`
	Namespace string   `json:"namespace"`
}

type FetchVectorsResponse struct {
	Vectors   []*FetchedVector `json:"vectors"`
	Namespace string           `json:"namespace"`
}

type FetchedVector struct {
	Id           string         `json:"id"`
	Values       []float32      `json:"values"`
	SparseValues *SparseValues  `json:"sparseValues"`
	Metadata     map[string]any `json:"metadata"`
}

func (ic *IndexClient) FetchVectors(ctx context.Context, params FetchVectorsParams) (*FetchVectorsResponse, error) {
	if params.Namespace == "" {
		return nil, fmt.Errorf("%w: namespace is required", ErrInvalidParams)
	}
	if len(params.Ids) == 0 {
		return nil, fmt.Errorf("%w: ids are required", ErrInvalidParams)
	}

	var respBody FetchVectorsResponse

	var urlParams string
	for _, id := range params.Ids {
		urlParams = strings.Join([]string{urlParams, fmt.Sprintf("&ids=%s", id)}, "")
	}
	if params.Namespace != "" {
		urlParams += fmt.Sprintf("&namespace=%s", params.Namespace)
	}
	resp, err := ic.reqClient.
		R().
		SetContentType("application/json").
		SetBody(params).
		SetSuccessResult(&respBody).
		SetContext(ctx).
		Get("/vectors/fetch?")
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
