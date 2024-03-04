package http

import (
	"encoding/json"
	"github-com/jairogloz/auth0-back/pkg/domain"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// GetTransactions returns the transactions
func GetTransactions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	txs := []domain.Transaction{
		{ID: "a"},
		{ID: "b"},
		{ID: "c"},
	}

	// Set content type as JSON for the response
	w.Header().Set("Content-Type", "application/json")

	// Encode and send back the same JSON body
	if err := json.NewEncoder(w).Encode(txs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
