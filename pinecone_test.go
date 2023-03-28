package pinecone

import (
	"os"
	"testing"
)

const (
	EnvTestPineconeAPIKey = "TEST_PINECONE_API_KEY"
)

func TestMain(m *testing.M) {
	if os.Getenv(EnvTestPineconeAPIKey) == "" {
		panic("TEST_PINECONE_API_KEY environment variable is not set")
	}

	os.Exit(m.Run())
}
