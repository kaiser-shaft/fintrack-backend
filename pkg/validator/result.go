package validator

import "net/http"

type Result struct {
	HasError bool
	Fields   []ValidationErrorResponse
}

func (res Result) StatusCode() int {
	if !res.HasError {
		return http.StatusOK
	}
	if len(res.Fields) > 0 {
		return http.StatusUnprocessableEntity
	}
	return http.StatusBadRequest
}
