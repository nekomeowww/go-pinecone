package pinecone

import (
	"os"
	"testing"
)

const (
	EnvTestPineconeAPIKey = "TEST_PINECONE_API_KEY"
)

func TestMain(m *testing.M) {
	os.Setenv(EnvTestPineconeAPIKey, "ad1793c0-da36-4f07-b345-887b73353fa5")

	if os.Getenv(EnvTestPineconeAPIKey) == "" {
		panic("TEST_PINECONE_API_KEY environment variable is not set")
	}

	os.Exit(m.Run())
}
