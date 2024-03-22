package handler

import (
	"encoding/json"
	"net/http"
)

func sendJSONResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	data, _ := json.Marshal(response)
	w.Write(data)
}

func methodNotAllowed(w http.ResponseWriter) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	e := newError("method not allowed", http.StatusMethodNotAllowed)
	http.Error(w, e.ToJson(), e.StatusCode)
}
