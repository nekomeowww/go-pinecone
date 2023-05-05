package pinecone

import "fmt"

// validateQueryParams validates the query parameters.
func validateQueryParams(params QueryParams) error {
	if params.TopK < 1 {
		return fmt.Errorf("%w: top k is required and must be greater than 0", ErrInvalidParams)
	}

	if params.Vector == nil && params.ID == "" {
		return fmt.Errorf("%w: vector or id is required", ErrInvalidParams)
	}

	if params.Vector != nil && params.ID != "" {
		return fmt.Errorf("%w: cannot specify both vector and id", ErrInvalidParams)
	}

	if params.SparseVector != nil && len(params.SparseVector.Values) != len(params.SparseVector.Indices) {
		return fmt.Errorf("%w: sparse vector values and indices must be the same length", ErrInvalidParams)
	}

	return nil
}

// validateDeleteVectorsParams validates the delete vectors parameters.
func validateDeleteVectorsParams(params DeleteVectorsParams) error {
	if (params.IDs == nil || len(params.IDs) == 0) && !params.DeleteAll {
		return fmt.Errorf("%w: ids or deleteAll is required", ErrInvalidParams)
	}

	if len(params.IDs) > 0 && params.DeleteAll {
		return fmt.Errorf("%w: cannot specify both ids and deleteAll", ErrInvalidParams)
	}

	return nil
}

// validateFetchVectorsParams validates the fetch vectors parameters.
func validateFetchVectorsParams(params FetchVectorsParams) error {
	if params.IDs == nil || len(params.IDs) < 1 {
		return fmt.Errorf("%w: ids is required", ErrInvalidParams)
	}

	return nil
}

// validateUpdateVectorParams validates the update vector parameters.
func validateUpdateVectorParams(params UpdateVectorParams) error {
	if params.ID == "" {
		return fmt.Errorf("%w: id is required", ErrInvalidParams)
	}

	if params.SparseValues != nil && len(params.SparseValues.Values) != len(params.SparseValues.Indices) {
		return fmt.Errorf("%w: sparse vector values and indices must be the same length", ErrInvalidParams)
	}

	return nil
}

// validateUpsertVectorsParams validates the upsert vectors parameters.
func validateUpsertVectorsParams(params UpsertVectorsParams) error {
	if params.Vectors == nil || len(params.Vectors) < 1 {
		return fmt.Errorf("%w: vectors is required", ErrInvalidParams)
	}

	for _, v := range params.Vectors {
		if v.SparseValues != nil && len(v.SparseValues.Values) != len(v.SparseValues.Indices) {
			return fmt.Errorf("%w: sparse vector values and indices must be the same length", ErrInvalidParams)
		}
	}

	return nil
}
