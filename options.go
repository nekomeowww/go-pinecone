package pinecone

type options struct {
	apiKey      string
	environment string
	projectName string
	indexName   string
}

type CallOptions struct {
	applyFunc func(o *options)
}

func applyCallOptions(callOptions []CallOptions, defaultOptions ...options) *options {
	o := new(options)
	if len(defaultOptions) > 0 {
		*o = defaultOptions[0]
	}

	for _, callOption := range callOptions {
		callOption.applyFunc(o)
	}

	return o
}

// WithAPIKey sets the API key to use for the call.
func WithAPIKey(apiKey string) CallOptions {
	return CallOptions{
		applyFunc: func(o *options) {
			o.apiKey = apiKey
		},
	}
}

// WithEnvironment sets the environment to use for the call.
func WithEnvironment(environment string) CallOptions {
	return CallOptions{
		applyFunc: func(o *options) {
			o.environment = environment
		},
	}
}

// WithProjectName sets the project name to use for the call.
func WithProjectName(projectName string) CallOptions {
	return CallOptions{
		applyFunc: func(o *options) {
			o.projectName = projectName
		},
	}
}

func WithIndexName(indexName string) CallOptions {
	return CallOptions{
		applyFunc: func(o *options) {
			o.indexName = indexName
		},
	}
}
