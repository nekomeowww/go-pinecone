package pinecone

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIndexClient_DeleteVectors(t *testing.T) {
	t.Run("should return error if namespace is empty", func(t *testing.T) {
		t.Run("Error", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			ic, err := NewIndexClient(
				WithAPIKey(""),
				WithEnvironment("us-central1-gcp"),
				WithIndexName("123"),
				WithProjectName("123"),
			)
			require.NoError(err)

			err = ic.DeleteVectors(context.Background(), DeleteVectorsParams{})
			require.Error(err)
			assert.ErrorIs(err, ErrRequestFailed)
			assert.EqualError(err, "API key is missing or invalid for the environment \"us-central1-gcp\". Check that the correct environment is specified.")
		})

	})
}

func TestIndexClient_FetchVectors(t *testing.T) {
	t.Run("should return error if namespace is empty", func(t *testing.T) {
		ic := IndexClient{}
		_, err := ic.FetchVectors(context.Background(), FetchVectorsParams{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestIndexClient_Query(t *testing.T) {
	t.Run("should return error if namespace is empty", func(t *testing.T) {
		ic := IndexClient{}
		_, err := ic.Query(context.Background(), QueryParams{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
