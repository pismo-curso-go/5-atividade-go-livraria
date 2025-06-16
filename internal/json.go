package internal

import (
	"encoding/json"
	"net/http"
)

type WriteJSONMessageParams struct {
	Writer     http.ResponseWriter
	HttpStatus int
	Data       any
}

type WriteJSONErrorParams struct {
	Writer     http.ResponseWriter
	HttpStatus int
	ErrMessage string
}

type dataWrapper struct {
	Data any `json:"data"`
}

func WriteJSONMessage(params WriteJSONMessageParams) {
	if params.HttpStatus == 0 {
		params.HttpStatus = http.StatusOK
	}

	message := dataWrapper{
		Data: params.Data,
	}

	params.Writer.Header().Set("Content-Type", "application/json")
	params.Writer.WriteHeader(params.HttpStatus)
	json.NewEncoder(params.Writer).Encode(message)
}

func WriteJSONError(params WriteJSONErrorParams) {
	if params.HttpStatus == 0 {
		params.HttpStatus = http.StatusInternalServerError
	}

	message := dataWrapper{
		Data: map[string]string{
			"error": params.ErrMessage,
		},
	}

	params.Writer.Header().Set("Content-Type", "application/json")
	params.Writer.WriteHeader(params.HttpStatus)
	json.NewEncoder(params.Writer).Encode(message)
}
