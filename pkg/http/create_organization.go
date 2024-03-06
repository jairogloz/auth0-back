package http

import (
	"context"
	"encoding/json"
	"github.com/auth0/go-auth0/management"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type OrganizationCreateParams struct {
	Name     *string `json:"name"`
	ClientId string  `json:"client_id"`
}

func (s *Server) CreateOrganization(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var organizationCreateParams OrganizationCreateParams

	// Decode the JSON body
	err := json.NewDecoder(r.Body).Decode(&organizationCreateParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newOrganization := &management.Organization{
		Name: organizationCreateParams.Name,
		Metadata: &map[string]string{
			"client_id": organizationCreateParams.ClientId,
		},
	}

	err = s.auth0API.Organization.Create(context.Background(), newOrganization)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content type as JSON for the response
	w.Header().Set("Content-Type", "application/json")

	// Encode and send back the same JSON body
	if err := json.NewEncoder(w).Encode(newOrganization); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func strPtr(s string) *string {
	return &s
}
