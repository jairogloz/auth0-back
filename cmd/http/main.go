package main

import (
	"context"
	"fmt"
	"github-com/jairogloz/auth0-back/pkg/auth0"
	pkgHttp "github-com/jairogloz/auth0-back/pkg/http"
	"github.com/auth0/go-auth0/management"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

// using https://github.com/auth0/go-auth0
func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading the .env file: %v", err)
	}

	// Get these from your Auth0 Application Dashboard.
	// The application needs to be a Machine To Machine authorized
	// to request access tokens for the Auth0 Management API,
	// with the desired permissions (scopes).
	domain := os.Getenv("AUTH0_DOMAIN")
	clientID := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")

	// Initialize a new client using a domain, client ID and client secret.
	// Alternatively you can specify an access token:
	// `management.WithStaticToken("token")`
	auth0API, err := management.New(
		domain,
		management.WithClientCredentials(context.TODO(), clientID, clientSecret), // Replace with a Context that better suits your usage
	)
	if err != nil {
		log.Fatalf("failed to initialize the auth0 management API client: %+v", err)
	}

	server := pkgHttp.NewServer(auth0API)

	router := httprouter.New()
	router.POST("/v1/public", server.Public)
	router.GET("/v1/transactions", auth0.EnsureValidTokenMiddleware(auth0.AuthroizeMddlw(server.GetTransactions)))
	router.POST("/v1/organizations", auth0.EnsureValidTokenMiddleware(auth0.AuthroizeMddlw(server.CreateOrganization)))

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
