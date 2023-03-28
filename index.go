package pinecone

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/pinecone-io/go-pinecone/pinecone_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// Index Pinecone Index client for vector operations
type Index struct {
	pinecone_grpc.VectorServiceClient
	conn *grpc.ClientConn
}

// Close closes the underlying gRPC connection.
func (i *Index) Close() {
	if i.conn != nil {
		i.conn.Close()
	}
}

// Index returns a new Index client.
//
// In best practice, index.Close() should be called
// after use since the index client is backed by a gRPC connection.
func (c *Client) Index(ctx context.Context, indexName string) (*Index, error) {
	tlsConfig := &tls.Config{}
	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", c.options.apiKey)
	pineconeTarget := fmt.Sprintf("%s-%s.svc.%s.pinecone.io:443", indexName, c.options.projectName, c.options.environment)
	conn, err := grpc.DialContext(
		ctx,
		pineconeTarget,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithAuthority(pineconeTarget),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &Index{
		VectorServiceClient: pinecone_grpc.NewVectorServiceClient(conn),
		conn:                conn,
	}, nil
}
