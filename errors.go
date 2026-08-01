package openapi

import "fmt"

// APIError 平台 OpenAPI 业务/HTTP 错误。
type APIError struct {
	Code       int
	Message    string
	RequestID  string
	HTTPStatus int
	RawBody    string
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Code != 0 {
		return fmt.Sprintf("openapi code=%d msg=%s requestId=%s", e.Code, e.Message, e.RequestID)
	}
	return fmt.Sprintf("openapi http=%d msg=%s", e.HTTPStatus, e.Message)
}

func AsAPIError(err error) (*APIError, bool) {
	if err == nil {
		return nil, false
	}
	e, ok := err.(*APIError)
	return e, ok
}
