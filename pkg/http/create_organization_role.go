package http

import (
	"encoding/json"
	"github-com/jairogloz/auth0-back/pkg/domain"
	"github.com/auth0/go-auth0/management"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type OrganizationRoleCreateParams struct {
	OrganizationId string `json:"organization_id"`
	Name           string `json:"name"`
}

func (s *Server) CreateOrganizationRole(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var organizationRoleCreateParams OrganizationRoleCreateParams

	// Decode the JSON body
	err := json.NewDecoder(r.Body).Decode(&organizationRoleCreateParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newRole := &management.Role{
		Name: domain.ToOrganizationRole(
			organizationRoleCreateParams.OrganizationId,
			organizationRoleCreateParams.Name),
	}

	err = s.auth0API.Role.Create(r.Context(), newRole)
	if err != nil {
		log.Println("error calling auth0API.Role.Create", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content type as JSON for the response
	w.Header().Set("Content-Type", "application/json")

	// Encode and send back the same JSON body
	if err := json.NewEncoder(w).Encode(newRole); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
