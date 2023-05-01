package pinecone

import (
	"bytes"
	"context"
	"fmt"
)

type Vector struct {
	Id       string            `json:"id"`
	Values   []float32         `json:"values"`
	Metadata map[string]string `json:"metadata"`
}

type UpsertVectorsParams struct {
	Vectors   []*Vector `json:"vectors"`
	Namespace string    `json:"namespace"`
}

type UpsertVectorsResponse struct {
	UpsertedCount int `json:"upsertedCount"`
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
