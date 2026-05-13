package render

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, message string, statusCode int, details any) {
	data := map[string]any{
		"error": message,
	}
	if details != nil {
		data["details"] = details
	}
	JSON(w, data, statusCode)
}
