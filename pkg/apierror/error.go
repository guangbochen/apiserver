package apierror

import (
	"errors"
	"fmt"

	"github.com/rancher/wrangler/v2/pkg/schemas/validation"
)

type APIError struct {
	Code      validation.ErrorCode
	Message   string
	Cause     error
	FieldName string
}

func NewAPIError(code validation.ErrorCode, message string) error {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

func NewFieldAPIError(code validation.ErrorCode, fieldName, message string) error {
	return &APIError{
		Code:      code,
		Message:   message,
		FieldName: fieldName,
	}
}

// WrapFieldAPIError will cause the API framework to log the underlying err before returning the APIError as a response.
// err WILL NOT be in the API response
func WrapFieldAPIError(err error, code validation.ErrorCode, fieldName, message string) error {
	return &APIError{
		Cause:     err,
		Code:      code,
		Message:   message,
		FieldName: fieldName,
	}
}

// WrapAPIError will cause the API framework to log the underlying err before returning the APIError as a response.
// err WILL NOT be in the API response
func WrapAPIError(err error, code validation.ErrorCode, message string) error {
	return &APIError{
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

func (a *APIError) Error() string {
	if a.FieldName != "" {
		return fmt.Sprintf("%s=%s: %s", a.FieldName, a.Code, a.Message)
	}
	return fmt.Sprintf("%s: %s", a.Code, a.Message)
}

func IsAPIError(err error) bool {
	var APIError *APIError
	ok := errors.As(err, &APIError)
	return ok
}

func IsConflict(err error) bool {
	var apiError *APIError
	if errors.As(err, &apiError) {
		return apiError.Code.Status == 409
	}

	return false
}
