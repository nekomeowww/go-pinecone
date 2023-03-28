package pinecone

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/samber/lo"
	"github.com/samber/mo"
)

// ListIndexes returns a list of your Pinecone
// indexes.
//
// API Reference: https://docs.pinecone.io/reference/list_indexes
func (c *Client) ListIndexes() ([]string, error) {
	var indexes []string
	resp, err := c.reqClient.
		R().
		SetSuccessResult(&indexes).
		Get("/databases")
	if err != nil {
		return make([]string, 0), err
	}
	if !resp.IsSuccessState() {
		buffer := new(bytes.Buffer)
		_, err := buffer.ReadFrom(resp.Body)
		if err != nil {
			return make([]string, 0), err
		}

		return make([]string, 0), fmt.Errorf("%w: %s, status code: %d", ErrRequestFailed, buffer.String(), resp.StatusCode)
	}

	return indexes, nil
}

type CreateIndexMetric string

const (
	CreateIndexMetricEuclidean  CreateIndexMetric = "euclidean"
	CreateIndexMetricCosine     CreateIndexMetric = "cosine"
	CreateIndexMetricDotProduct CreateIndexMetric = "dotproduct"
)

type CreateIndexPodType string

const (
	CreateIndexPodTypeS1 CreateIndexPodType = "s1"
	CreateIndexPodTypeP1 CreateIndexPodType = "p1"
	CreateIndexPodTypeP2 CreateIndexPodType = "p2"
)

type CreateIndexPodSize string

const (
	CreateIndexPodSize1 CreateIndexPodSize = "1"
	CreateIndexPodSize2 CreateIndexPodSize = "2"
	CreateIndexPodSize4 CreateIndexPodSize = "4"
	CreateIndexPodSize8 CreateIndexPodSize = "8"
)

type CreateIndexParams struct {
	// Required. The name of the index to be created.
	// The maximum length is 45 characters.
	Name string
	// Required. The dimensions of the vectors to be
	// inserted in the index
	Dimension int
	// The distance metric to be used for similarity
	// search. You can use 'euclidean', 'cosine', or
	// 'dotproduct'.
	Metric mo.Option[CreateIndexMetric]
	// The number of pods for the index to use,
	// including replicas.
	Pods mo.Option[int]
	// The number of replicas. Replicas duplicate
	// your index. They provide higher availability
	// and throughput.
	Replicas mo.Option[int]
	// The type of pod to use. One of s1, p1, or p2
	PodType mo.Option[CreateIndexPodType]
	PodSize mo.Option[CreateIndexPodSize]
	// Configuration for the behavior of Pinecone's
	// internal metadata index. By default, all
	// metadata is indexed; when metadata_config is
	// present, only specified metadata fields are
	// indexed. To specify metadata fields to index,
	// provide a JSON object of the following form:
	//	{"indexed": ["example_metadata_field"]}
	MetadataConfig mo.Option[map[string]string]
	// The name of the collection to create an index
	// from
	SourceCollection mo.Option[string]
}

type CreateIndexBodyParams struct {
	Name             string             `json:"name"`
	Dimension        int                `json:"dimension"`
	Metric           *CreateIndexMetric `json:"metric,omitempty"`
	Pods             *int               `json:"pods,omitempty"`
	Replicas         *int               `json:"replicas,omitempty"`
	PodType          *string            `json:"pod_type,omitempty"`
	MetadataConfig   map[string]string  `json:"metadata_config,omitempty"`
	SourceCollection *string            `json:"source_collection,omitempty"`
}

// CreateIndex creates a Pinecone index. You can use
// it to specify the measure of similarity, the
// dimension of vectors to be stored in the index,
// the numbers of replicas to use, and more.
//
// API Reference: https://docs.pinecone.io/reference/create_index
func (c *Client) CreateIndex(ctx context.Context, params CreateIndexParams) error {
	if params.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidParams)
	}
	if params.Dimension <= 0 {
		return fmt.Errorf("%w: dimension is required", ErrInvalidParams)
	}
	if params.PodSize.IsPresent() && params.PodType.IsAbsent() {
		return fmt.Errorf("%w: pod_type is required when pod_size is specified", ErrInvalidParams)
	}
	if params.PodType.IsPresent() && params.PodSize.IsAbsent() {
		return fmt.Errorf("%w: pod_size is required when pod_type is specified", ErrInvalidParams)
	}

	var body CreateIndexBodyParams
	body.Name = params.Name
	body.Dimension = params.Dimension
	if params.Metric.IsPresent() {
		body.Metric = lo.ToPtr(params.Metric.MustGet())
	}
	if params.Pods.IsPresent() {
		body.Pods = lo.ToPtr(params.Pods.MustGet())
	}
	if params.Replicas.IsPresent() {
		body.Replicas = lo.ToPtr(params.Replicas.MustGet())
	}
	if params.PodType.IsPresent() && params.PodSize.IsPresent() {
		podType := params.PodType.MustGet()
		podSize := params.PodSize.MustGet()
		podTypeSize := string(podType) + "." + string(podSize)
		body.PodType = lo.ToPtr(podTypeSize)
	}
	if params.MetadataConfig.IsPresent() {
		body.MetadataConfig = params.MetadataConfig.MustGet()
	}
	if params.SourceCollection.IsPresent() {
		body.SourceCollection = lo.ToPtr(params.SourceCollection.MustGet())
	}

	resp, err := c.reqClient.
		R().
		SetContentType("application/json").
		SetBody(body).
		SetContext(ctx).
		Post("/databases")
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

type DescribeIndexResponse struct {
	Database Database `json:"database"`
	Status   Status   `json:"status"`
}

type Database struct {
	Name      string `json:"name"`
	Metric    string `json:"metric"`
	Dimension int    `json:"dimension"`
	Replicas  int    `json:"replicas"`
	Shards    int    `json:"shards"`
	Pods      int    `json:"pods"`
	PodType   string `json:"pod_type"`
}

type Status struct {
	Waiting []interface{} `json:"waiting"`
	Crashed []interface{} `json:"crashed"`
	Host    string        `json:"host"`
	Port    int           `json:"port"`
	State   string        `json:"state"`
	Ready   bool          `json:"ready"`
}

// DescribeIndex gets a description of an index.
//
// API Reference: https://docs.pinecone.io/reference/describe_index
func (c *Client) DescribeIndex(ctx context.Context, indexName string) (*DescribeIndexResponse, error) {
	if indexName == "" {
		return nil, fmt.Errorf("%w: index name is required", ErrInvalidParams)
	}

	var respBody DescribeIndexResponse
	resp, err := c.reqClient.
		R().
		SetContentType("application/json").
		SetSuccessResult(&respBody).
		SetContext(ctx).
		Get("/databases/" + indexName)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() && resp.StatusCode == http.StatusNotFound {
		return nil, ErrIndexNotFound
	} else if !resp.IsSuccessState() {
		buffer := new(bytes.Buffer)
		_, err := buffer.ReadFrom(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%w: %s, status code: %d", ErrRequestFailed, buffer.String(), resp.StatusCode)
	}

	return &respBody, nil
}

// DeleteIndex deletes an existing index.
//
// API Reference: https://docs.pinecone.io/reference/delete_index
func (c *Client) DeleteIndex(ctx context.Context, indexName string) error {
	if indexName == "" {
		return fmt.Errorf("%w: index name is required", ErrInvalidParams)
	}

	resp, err := c.reqClient.
		R().
		SetContext(ctx).
		Delete("/databases/" + indexName)
	if err != nil {
		return err
	}
	if !resp.IsSuccessState() && resp.StatusCode == http.StatusNotFound {
		return ErrIndexNotFound
	} else if !resp.IsSuccessState() {
		buffer := new(bytes.Buffer)
		_, err := buffer.ReadFrom(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("%w: %s, status code: %d", ErrRequestFailed, buffer.String(), resp.StatusCode)
	}

	return nil
}

type ConfigureIndexParams struct {
	// The name of the index
	IndexName string
	// The number of replicas. Replicas duplicate
	// your index. They provide higher availability
	// and throughput.
	Replicas mo.Option[int]
	// The type of pod to use. One of s1, p1, or p2
	PodType mo.Option[CreateIndexPodType]
	PodSize mo.Option[CreateIndexPodSize]
}

type ConfigureIndexBodyParams struct {
	Replicas *int    `json:"replicas,omitempty"`
	PodType  *string `json:"pod_type,omitempty"`
}

// ConfigureIndex specifies the pod type and number of replicas for an index.
//
// API Reference: https://docs.pinecone.io/reference/configure_index
func (c *Client) ConfigureIndex(ctx context.Context, params ConfigureIndexParams) error {
	if params.IndexName == "" {
		return fmt.Errorf("%w: index name is required", ErrInvalidParams)
	}
	if params.PodSize.IsAbsent() && params.PodType.IsAbsent() && params.Replicas.IsAbsent() {
		return fmt.Errorf("%w: at least one of replicas, pod_type or pod_size is required", ErrInvalidParams)
	}
	if params.PodSize.IsPresent() && params.PodType.IsAbsent() {
		return fmt.Errorf("%w: pod_type is required when pod_size is specified", ErrInvalidParams)
	}
	if params.PodType.IsPresent() && params.PodSize.IsAbsent() {
		return fmt.Errorf("%w: pod_size is required when pod_type is specified", ErrInvalidParams)
	}

	resp, err := c.reqClient.
		R().
		SetContentType("application/json").
		SetContext(ctx).
		Patch("/databases/" + params.IndexName)
	if err != nil {
		return err
	}
	if !resp.IsSuccessState() && resp.StatusCode == http.StatusNotFound {
		return ErrIndexNotFound
	} else if !resp.IsSuccessState() {
		buffer := new(bytes.Buffer)
		_, err := buffer.ReadFrom(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("%w: %s, status code: %d", ErrRequestFailed, buffer.String(), resp.StatusCode)
	}

	return nil
}
