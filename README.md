# go-pinecone

[![Go Reference](https://pkg.go.dev/badge/badge/github.com/nekomeowww/go-pinecone.svg)](https://pkg.go.dev/badge/github.com/nekomeowww/go-pinecone)
[![Go Report](https://goreportcard.com/badge/github.com/nekomeowww/go-pinecone)](https://goreportcard.com/report/github.com/nekomeowww/go-pinecone)
[![Testing](https://github.com/nekomeowww/go-pinecone/actions/workflows/ci.yml/badge.svg)](https://github.com/nekomeowww/go-pinecone/actions/workflows/ci.yml)
[![Building](https://github.com/nekomeowww/go-pinecone/actions/workflows/build.yml/badge.svg)](https://github.com/nekomeowww/go-pinecone/actions/workflows/build.yml)

---

This package is an API client for Pinecone, a SaaS vector database. With this package, users can easily perform Index and Vector data operations in Golang projects.

## Introduction

Pinecone is a cloud-based vector database that enables users to store large-scale vector data and query them efficiently.

This repo aims to provide a simple and easy-to-use client for Golang users to interface with Pinecone. It supports all the operations that are available through the Pinecone API, such as creating and deleting indexes, inserting and querying vectors, and modifying index metadata, among other functionalities.

## Installation

To install, simply run the following command:

```sh
go get -u github.com/nekomeowww/go-pinecone
```

## Documentation

For a complete reference of the functions and types, please refer to the [godoc documentation](https://pkg.go.dev/github.com/nekomeowww/go-pinecone).

## Get started

### Initialize a new Pinecone client

```go
package main

import (
    pinecone "github.com/nekomeowww/go-pinecone"
)

func main() {
    p, err := pinecone.New(
        pinecone.WithAPIKey("YOUR_API_KEY"),
        pinecone.WithEnvironment("YOUR_ACCOUNT_REGION"),
        pinecone.WithProjectName("YOUR_PROJECT_NAME"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Do something with the client
}
```

### Enables debug and HTTP request dump

```go
package main

import (
    pinecone "github.com/nekomeowww/go-pinecone"
)

func main() {
    p, err := pinecone.New(
        pinecone.WithAPIKey("YOUR_API_KEY"),
        pinecone.WithEnvironment("YOUR_ACCOUNT_REGION"),
        pinecone.WithProjectName("YOUR_PROJECT_NAME"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Enable debug
    p = p.Debug()
}
```

### Establish a connection to interact with Vectors

```go
package main

import (
    pinecone "github.com/nekomeowww/go-pinecone"
)

func main() {
    p, err := pinecone.New(
        pinecone.WithAPIKey("YOUR_API_KEY"),
        pinecone.WithEnvironment("YOUR_ACCOUNT_REGION"),
        pinecone.WithProjectName("YOUR_PROJECT_NAME"),
    )
    if err != nil {
        log.Fatal(err)
    }
}
```


## Initialize a new Index client

Vector operations are performed on a given index, we initialize a new index client like:

```go
package main

import (
	"context"
	"fmt"
	pinecone "github.com/nekomeowww/go-pinecone"
	"log"
)

func main() {
	client, err := pinecone.NewIndexClient(
		pinecone.WithIndexName("YOUR_INDEX_NAME"),
		pinecone.WithEnvironment("YOUR_ACCOUNT_REGION"),
		pinecone.WithProjectName("YOUR_PROJECT_NAME"),
		pinecone.WithAPIKey("YOUR_API_KEY"),
	)

	if err != nil {
		log.Fatal(err)
	}
}
```

To describe index stats we can use:

```go
	ctx := context.Background()

	params := pinecone.DescribeIndexStatsParams{}
	resp, err := client.DescribeIndexStats(ctx, params)
	if err != nil {
		panic(err)
	}
	for k, v := range resp.Namespaces {
		fmt.Printf("%s: %+v\n", k, v)
	}
```

To Upsert Vectors we do such as:
    
```go
    // Use your own embedding client/library such as OpenAI embedding API
    myVectorEmbedding := embedding.New("this is a sentence I want to embed")
	
    ctx := context.Background()
    params := pinecone.UpsertVectorsParams{
        Vectors: []*pinecone.Vector{
            {
                ID: "YOUR_VECTOR_ID",
                Values: myVectorEmbedding,
                Metadata: map[string]any{"foo": "bar"}
            },
        },
    }
    resp, err := client.UpsertVectors(ctx, params)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", resp)
```

For a complete reference of the functions and types, please refer to the [godoc documentation](https://pkg.go.dev/github.com/nekomeowww/go-pinecone).

## Contributing

We welcome contributions from the Golang community! If you'd like to contribute, please follow these steps:

1. Fork this repository
2. Create a new branch for your changes
3. Make the changes
4. Commit your changes with a meaningful commit message
5. Create a pull request

## Acknowledgements

- Official Pinecone Index Client [go-pinecone](https://github.com/pinecone-io/go-pinecone)

## License

Released under the [MIT License](LICENSE).
