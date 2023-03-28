package pinecone

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexOperations(t *testing.T) {
	t.Run("CreateIndex", func(t *testing.T) {
		t.Run("Error", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			c, err := New(
				WithAPIKey(""),
				WithEnvironment("us-central1-gcp"),
			)
			require.NoError(err)

			err = c.CreateIndex(context.Background(), CreateIndexParams{
				Name:      "test-index",
				Dimension: 10,
			})
			require.Error(err)
			assert.ErrorIs(err, ErrRequestFailed)
			assert.EqualError(fmt.Errorf("%w: %s, status code: %d", ErrRequestFailed, `API key is missing or invalid for the environment "us-central1-gcp". Check that the correct environment is specified.`, 401), err.Error())
		})

		t.Run("Success", func(t *testing.T) {
			require := require.New(t)

			c, err := New(
				WithAPIKey(os.Getenv(EnvTestPineconeAPIKey)),
				WithEnvironment("us-central1-gcp"),
			)
			require.NoError(err)

			err = c.CreateIndex(context.Background(), CreateIndexParams{
				Name:      "test-index",
				Dimension: 10,
			})
			require.NoError(err)
		})
	})

	time.Sleep(time.Minute)

	t.Run("ListIndexes", func(t *testing.T) {
		t.Run("Error", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			c, err := New(
				WithAPIKey(""),
				WithEnvironment("us-central1-gcp"),
			)
			require.NoError(err)

			indexes, err := c.ListIndexes()
			require.Error(err)
			require.Empty(indexes)
			assert.ErrorIs(err, ErrRequestFailed)
			assert.EqualError(fmt.Errorf("%w: %s, status code: %d", ErrRequestFailed, `API key is missing or invalid for the environment "us-central1-gcp". Check that the correct environment is specified.`, 401), err.Error())
		})

		t.Run("Success", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			c, err := New(
				WithAPIKey(os.Getenv(EnvTestPineconeAPIKey)),
				WithEnvironment("us-central1-gcp"),
			)
			require.NoError(err)

			indexes, err := c.ListIndexes()
			require.NoError(err)
			assert.NotNil(indexes)
			assert.Len(indexes, 1)
		})
	})

	t.Run("DescribeIndex", func(t *testing.T) {
		t.Run("Error", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			c, err := New(
				WithAPIKey(""),
				WithEnvironment("us-central1-gcp"),
			)
			require.NoError(err)

			resp, err := c.DescribeIndex(context.Background(), "test-index")
			require.Error(err)
			require.Nil(resp)
			assert.ErrorIs(err, ErrRequestFailed)
			assert.EqualError(fmt.Errorf("%w: %s, status code: %d", ErrRequestFailed, `API key is missing or invalid for the environment "us-central1-gcp". Check that the correct environment is specified.`, 401), err.Error())
		})

		t.Run("NotFound", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			c, err := New(
				WithAPIKey(os.Getenv(EnvTestPineconeAPIKey)),
				WithEnvironment("us-central1-gcp"),
			)
			require.NoError(err)

			resp, err := c.DescribeIndex(context.Background(), fmt.Sprintf("test-index-%d", time.Now().UnixMilli()))
			require.Error(err)
			assert.ErrorIs(err, ErrIndexNotFound)
			assert.Nil(resp)
		})

		t.Run("Success", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			c, err := New(
				WithAPIKey(os.Getenv(EnvTestPineconeAPIKey)),
				WithEnvironment("us-central1-gcp"),
			)
			require.NoError(err)

			resp, err := c.DescribeIndex(context.Background(), "test-index")
			require.NoError(err)
			require.NotNil(resp)
			assert.Equal("test-index", resp.Database.Name)
		})
	})

	t.Run("ConfigureIndex", func(t *testing.T) {
		t.Run("Error", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			c, err := New(
				WithAPIKey(""),
				WithEnvironment("us-central1-gcp"),
			)
			require.NoError(err)

			err = c.ConfigureIndex(context.Background(), ConfigureIndexParams{
				IndexName: "test-index",
				Replicas:  mo.Some(2),
			})
			require.Error(err)
			assert.ErrorIs(err, ErrRequestFailed)
			assert.EqualError(fmt.Errorf("%w: %s, status code: %d", ErrRequestFailed, `API key is missing or invalid for the environment "us-central1-gcp". Check that the correct environment is specified.`, 401), err.Error())
		})

		t.Run("NotFound", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			c, err := New(
				WithAPIKey(os.Getenv(EnvTestPineconeAPIKey)),
				WithEnvironment("us-central1-gcp"),
			)
			require.NoError(err)

			err = c.ConfigureIndex(context.Background(), ConfigureIndexParams{
				IndexName: fmt.Sprintf("test-index-%d", time.Now().UnixMilli()),
				Replicas:  mo.Some(2),
			})
			require.Error(err)
			assert.ErrorIs(err, ErrIndexNotFound)
		})

		t.Run("Success", func(t *testing.T) {
			require := require.New(t)

			c, err := New(
				WithAPIKey(os.Getenv(EnvTestPineconeAPIKey)),
				WithEnvironment("us-central1-gcp"),
			)
			require.NoError(err)

			err = c.ConfigureIndex(context.Background(), ConfigureIndexParams{
				IndexName: "test-index",
				Replicas:  mo.Some(2),
			})
			require.NoError(err)
		})
	})

	t.Run("DeleteIndex", func(t *testing.T) {
		t.Run("Error", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			c, err := New(
				WithAPIKey(""),
				WithEnvironment("us-central1-gcp"),
			)
			require.NoError(err)

			err = c.DeleteIndex(context.Background(), "test-index")
			require.Error(err)
			assert.ErrorIs(err, ErrRequestFailed)
			assert.EqualError(fmt.Errorf("%w: %s, status code: %d", ErrRequestFailed, `API key is missing or invalid for the environment "us-central1-gcp". Check that the correct environment is specified.`, 401), err.Error())
		})

		t.Run("NotFound", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			c, err := New(
				WithAPIKey(os.Getenv(EnvTestPineconeAPIKey)),
				WithEnvironment("us-central1-gcp"),
			)
			require.NoError(err)

			err = c.DeleteIndex(context.Background(), fmt.Sprintf("test-index-%d", time.Now().UnixMilli()))
			require.Error(err)
			assert.ErrorIs(err, ErrIndexNotFound)
		})

		t.Run("Success", func(t *testing.T) {
			require := require.New(t)

			c, err := New(
				WithAPIKey(os.Getenv(EnvTestPineconeAPIKey)),
				WithEnvironment("us-central1-gcp"),
			)
			require.NoError(err)

			err = c.DeleteIndex(context.Background(), "test-index")
			require.NoError(err)
		})
	})
}
