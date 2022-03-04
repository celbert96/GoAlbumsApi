package models

import (
	"fmt"
	"net/http"
)

type ErrResponse struct {
	ErrorMessage string `json:"error_message"`
}

func ErrResponseForHttpStatus(status int) ErrResponse {
	switch status {
	case http.StatusForbidden:
		return ErrResponse{ErrorMessage: "access denied"}
	case http.StatusNotFound:
		return ErrResponse{ErrorMessage: "resource not found"}
	case http.StatusBadRequest:
		return ErrResponse{ErrorMessage: "bad request"}
	case http.StatusInternalServerError:
		return ErrResponse{ErrorMessage: "internal server error"}
	default:
		return ErrResponse{ErrorMessage: fmt.Sprintf("unknown error: code %d", status)}
	}
}
