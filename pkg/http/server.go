package http

import "github.com/auth0/go-auth0/management"

type Server struct {
	auth0API *management.Management
}

func NewServer(a *management.Management) *Server {
	return &Server{auth0API: a}
}
