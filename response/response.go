package response

import (
	"encoding/json"
	"net/http"
)

func JSONSuccess(w http.ResponseWriter, data any, pagination any, meta any) {
	response := map[string]any{
		"message": "Success",
		"data":    data,
	}

	if pagination != nil {
		response["pagination"] = pagination
	}

	if meta != nil {
		response["meta"] = meta
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func JSONError(w http.ResponseWriter, status int, message string, errorDetails map[string][]string) {
	response := map[string]interface{}{
		"message": message,
		"error":   errorDetails,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
