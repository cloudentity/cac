package storage

func NewWithID[T any](id string, data T) *WithID[T] {
	return &WithID[T]{
		ID:    id,
		Other: data,
	}
}

type WithID[T any] struct {
	ID string `json:"id" yaml:"id"`
	// nolint
	Other T `json:",inline,squash" yaml:",inline"`
}
