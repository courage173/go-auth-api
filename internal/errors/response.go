package errors

import (
	"net/http"
	"sort"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)


type ErrorResponse struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Errors interface{} `json:"errors,omitempty"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

func (e ErrorResponse) StatusCode() int {
	return e.Status
}

func InternalServerError(msg string) ErrorResponse {
	if msg == ""{
		msg = "An error occured while processing the request"
	}
	
	return ErrorResponse{
        Status:  http.StatusInternalServerError,
        Message: "Internal Server Error",
    }
}

func NotFound(msg string) ErrorResponse {
	if msg == ""{
        msg = "The resource was not Found"
    }
    
    return ErrorResponse{
        Status:  http.StatusNotFound,
        Message: msg,
    }
}

func Unauthorized(msg string) ErrorResponse {
	if msg == "" {
		msg = "You are not authenticated to perform the requested action."
	}
	return ErrorResponse{
		Status:  http.StatusUnauthorized,
		Message: msg,
	}
}


func Forbidden(msg string) ErrorResponse {
	if msg == "" {
		msg = "You are not authorized to perform the requested action."
	}
	return ErrorResponse{
		Status:  http.StatusForbidden,
		Message: msg,
	}
}

// BadRequest creates a new error response representing a bad request (HTTP 400)
func BadRequest(msg string) ErrorResponse {
	if msg == "" {
		msg = "Your request is in a bad format."
	}
	return ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: msg,
	}
}

type invalidField struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// InvalidInput creates a new error response representing a data validation error (HTTP 400).
func InvalidInput(errs validation.Errors) ErrorResponse {
	var details []invalidField
	var fields []string
	for field := range errs {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	for _, field := range fields {
		details = append(details, invalidField{
			Field: field,
			Error: errs[field].Error(),
		})
	}

	return ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: "There is some problem with the data you submitted.",
		Errors: details,
	}
}