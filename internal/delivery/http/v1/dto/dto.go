package dto

type ListResponse[T any] struct {
	Data []T `json:"data"`
}

func MapSlice[S any, R any](items []S, mapper func(S) R) []R {
	res := make([]R, len(items))
	for i, item := range items {
		res[i] = mapper(item)
	}
	return res
}
