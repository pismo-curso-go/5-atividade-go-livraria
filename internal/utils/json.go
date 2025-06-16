package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func WriteJSONError(w http.ResponseWriter, err *APIError) {
	WriteJSONResponse(w, err.Code, map[string]interface{}{
		"error": err.Message,
		"code":  err.Code,
	})
}

func DecodeJSONRequest(r *http.Request, dest interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dest); err != nil {
		return ErrInvalidJSON
	}

	return nil
}
