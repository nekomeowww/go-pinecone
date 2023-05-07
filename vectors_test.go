package pinecone

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidators(t *testing.T) {

	t.Run("validate query params", func(t *testing.T) {
		t.Run("should return error if top k is less than 1", func(t *testing.T) {
			params := QueryParams{
				TopK: 0,
			}

			err := validateQueryParams(params)

			require.Error(t, err)
			if err.Error() != "invalid params: top k is required and must be greater than 0" {
				t.Errorf("expected: invalid params: top k is required and must be greater than 0, got: %s", err.Error())
			}
		})

		t.Run("should not return error if top k is 1 and id is specified", func(t *testing.T) {
			params := QueryParams{
				TopK: 1,
				ID:   "id",
			}

			err := validateQueryParams(params)
			require.NoError(t, err)
		})

		t.Run("should return error if vector and id are both empty", func(t *testing.T) {
			params := QueryParams{
				TopK: 1,
			}

			err := validateQueryParams(params)
			require.Error(t, err)
			if err.Error() != "invalid params: vector or id is required" {
				t.Errorf("expected: invalid params: vector or id is required, got: %s", err.Error())
			}
		})

		t.Run("should return error if vector and id are both specified", func(t *testing.T) {
			params := QueryParams{
				TopK:   1,
				Vector: []float32{1.0, 2.0},
				ID:     "id",
			}

			err := validateQueryParams(params)
			require.Error(t, err)
			if err.Error() != "invalid params: cannot specify both vector and id" {
				t.Errorf("expected: invalid params: cannot specify both vector and id, got: %s", err.Error())
			}
		})

		t.Run("should not return error if vector is specified and id is empty", func(t *testing.T) {
			params := QueryParams{
				TopK:   1,
				Vector: []float32{1.0, 2.0},
			}

			err := validateQueryParams(params)
			require.NoError(t, err)
		})

		t.Run("should return error if sparse vector values and indices are not the same length", func(t *testing.T) {
			params := QueryParams{
				TopK:         1,
				ID:           "id",
				SparseVector: &SparseVector{Values: []float32{1.0}, Indices: []int32{1, 2}},
			}

			err := validateQueryParams(params)
			require.Error(t, err)
			if err.Error() != "invalid params: sparse vector values and indices must be the same length" {
				t.Errorf("expected: invalid params: sparse vector values and indices must be the same length, got: %s", err.Error())
			}
		})
	})

	t.Run("validate delete vectors params", func(t *testing.T) {
		t.Run("should return error if IDs and DeleteAll are nil", func(t *testing.T) {
			params := DeleteVectorsParams{}

			err := validateDeleteVectorsParams(params)
			require.Error(t, err)
			if err.Error() != "invalid params: ids or deleteAll is required" {
				t.Errorf("expected: invalid params: ids or deleteAll is required, got: %s", err.Error())
			}
		})

		t.Run("should return error if IDs and DeleteAll are both specified", func(t *testing.T) {
			params := DeleteVectorsParams{
				IDs:       []string{"id"},
				DeleteAll: true,
			}

			err := validateDeleteVectorsParams(params)
			require.Error(t, err)
			if err.Error() != "invalid params: cannot specify both ids and deleteAll" {
				t.Errorf("expected: invalid params: cannot specify both ids and deleteAll, got: %s", err.Error())
			}
		})

		t.Run("should return error if IDs is empty", func(t *testing.T) {
			params := DeleteVectorsParams{
				IDs: []string{},
			}

			err := validateDeleteVectorsParams(params)
			require.Error(t, err)
			if err.Error() != "invalid params: ids or deleteAll is required" {
				t.Errorf("expected: invalid params: ids must not be empty, got: %s", err.Error())
			}
		})

		t.Run("should not return error if DeleteAll is true", func(t *testing.T) {
			params := DeleteVectorsParams{
				DeleteAll: true,
			}

			err := validateDeleteVectorsParams(params)
			require.NoError(t, err)
		})

		t.Run("should not return error if IDs is not empty", func(t *testing.T) {
			params := DeleteVectorsParams{
				IDs: []string{"id"},
			}

			err := validateDeleteVectorsParams(params)
			require.NoError(t, err)
		})
	})

	t.Run("validate fetch vectors params", func(t *testing.T) {
		t.Run("should return error if ids is nil", func(t *testing.T) {
			params := FetchVectorsParams{
				IDs: nil,
			}

			err := validateFetchVectorsParams(params)
			require.Error(t, err)
			if err.Error() != "invalid params: ids is required" {
				t.Errorf("expected: invalid params: ids is required, got: %s", err.Error())
			}
		})

		t.Run("should return error if ids is empty", func(t *testing.T) {
			params := FetchVectorsParams{
				IDs: []string{},
			}

			err := validateFetchVectorsParams(params)
			require.Error(t, err)
			if err.Error() != "invalid params: ids is required" {
				t.Errorf("expected: invalid params: ids must not be empty, got: %s", err.Error())
			}
		})

		t.Run("should not return error if ids is not empty", func(t *testing.T) {
			params := FetchVectorsParams{
				IDs: []string{"id"},
			}

			err := validateFetchVectorsParams(params)
			require.NoError(t, err)
		})
	})

	t.Run("validate update vector params", func(t *testing.T) {
		t.Run("should return error if id is empty", func(t *testing.T) {
			params := UpdateVectorParams{
				ID: "",
			}

			err := validateUpdateVectorParams(params)
			require.Error(t, err)
			if err.Error() != "invalid params: id is required" {
				t.Errorf("expected: invalid params: id is required, got: %s", err.Error())
			}
		})

		t.Run("should return error if params SparseValues and SparseIndices are not the same length", func(t *testing.T) {
			params := UpdateVectorParams{
				ID: "id",
				SparseValues: &SparseVector{
					Values:  []float32{1.0},
					Indices: []int32{1, 2},
				},
			}

			err := validateUpdateVectorParams(params)
			require.Error(t, err)
			if err.Error() != "invalid params: sparse vector values and indices must be the same length" {
				t.Errorf("expected: invalid params: sparse vector values and indices must be the same length, got: %s", err.Error())
			}
		})

		t.Run("should not return error if params SparseValues and SparseIndices are the same length", func(t *testing.T) {
			params := UpdateVectorParams{
				ID: "id",
				SparseValues: &SparseVector{
					Values:  []float32{1.0, 2.0},
					Indices: []int32{1, 2},
				},
			}

			err := validateUpdateVectorParams(params)
			require.NoError(t, err)
		})

	})

	t.Run("validate upsert vectors params", func(t *testing.T) {
		t.Run("should return error if vectors is nil", func(t *testing.T) {
			params := UpsertVectorsParams{
				Vectors: nil,
			}
			err := validateUpsertVectorsParams(params)
			require.Error(t, err)
			if err.Error() != "invalid params: vectors is required" {
				t.Errorf("expected: invalid params: vectors is required, got: %s", err.Error())
			}
		})

		t.Run("should return error if vectors is empty", func(t *testing.T) {
			params := UpsertVectorsParams{
				Vectors: []*Vector{},
			}
			err := validateUpsertVectorsParams(params)
			require.Error(t, err)
			if err.Error() != "invalid params: vectors is required" {
				t.Errorf("expected: invalid params: vectors is required, got: %s", err.Error())
			}
		})

		t.Run("should return error if sparse vectors values and indices are not the same length", func(t *testing.T) {
			params := UpsertVectorsParams{
				Vectors: []*Vector{
					{
						ID: "id",
						SparseValues: &SparseVector{
							Values:  []float32{1.0},
							Indices: []int32{1, 2},
						},
					},
				},
			}
			err := validateUpsertVectorsParams(params)
			require.Error(t, err)
			if err.Error() != "invalid params: sparse vector values and indices must be the same length" {
				t.Errorf("expected: invalid params: sparse vector values and indices must be the same length, got: %s", err.Error())
			}
		})

		t.Run("should not return error if sparse vectors values and indices are the same length", func(t *testing.T) {
			params := UpsertVectorsParams{
				Vectors: []*Vector{
					{
						ID: "id",
						SparseValues: &SparseVector{
							Values:  []float32{1.0, 2.0},
							Indices: []int32{1, 2},
						},
					},
				},
			}
			err := validateUpsertVectorsParams(params)
			require.NoError(t, err)
		})

		t.Run("should not return error if sparse vectors values and indices are nil", func(t *testing.T) {
			params := UpsertVectorsParams{
				Vectors: []*Vector{
					{
						ID: "id",
						SparseValues: &SparseVector{
							Values:  nil,
							Indices: nil,
						},
					},
				},
			}
			err := validateUpsertVectorsParams(params)
			require.NoError(t, err)
		})
	})
}
