package http

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Public is an http handler that echoes the request body.
func (s *Server) Public(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var requestBody map[string]interface{}

	// Decode the JSON body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set content type as JSON for the response
	w.Header().Set("Content-Type", "application/json")

	// Encode and send back the same JSON body
	if err := json.NewEncoder(w).Encode(requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
