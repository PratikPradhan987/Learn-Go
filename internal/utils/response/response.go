package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error string `json:"error"`
}

const (
	StatusOK   = 200
	StatusError = "Error"
    StatusBadRequest = 400
    StatusInternalServerError = 500
    NotFound          = 404
    InternalError     = 500
    Error             = "error"
    StatusText       = "status"
    ErrorText         = "error_message"
    SuccessText      = "success"
    MessageText      = "message"
    DataText         = "data"
    StatusOk         = "ok"
    StatusFail       = "fail"
    MessageOk        = "Request processed successfully"
    MessageFail      = "An error occurred while processing the request"
    MessageNotFound  = "Requested resource not found"
    MessageInternal  = "Internal server error"
    MessageError     = "An error occurred"
    DataKey          = "data"
    StatusKey        = "status"
    ErrorKey         = "error"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
        Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors)  Response{
	var errMsgs []string

	for _, err := range errs{
		switch err.ActualTag() {
			case "required":
                errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
            case "min":
                errMsgs = append(errMsgs, fmt.Sprintf("field %s is min", err.Field()))
            case "max":
                errMsgs = append(errMsgs, fmt.Sprintf("field %s is max", err.Field()))
            default:
                errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
        Error:  strings.Join(errMsgs, ", "),
	}
}