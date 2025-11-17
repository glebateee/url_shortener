package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	statusOK    = "OK"
	statusError = "Error"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func OK() Response {
	return Response{
		Status: statusOK,
	}
}

func Error(err string) Response {
	return Response{
		Status: statusError,
		Error:  err,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid url", err.Field()))
		case "default":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}
	return Response{
		Status: statusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
