package main

import (
	"fmt"
	"github-com/jairogloz/auth0-back/pkg/auth0"
	pkgHttp "github-com/jairogloz/auth0-back/pkg/http"
	"github.com/joho/godotenv"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading the .env file: %v", err)
	}

	router := httprouter.New()
	router.POST("/v1/public", pkgHttp.Public)
	router.GET("/v1/transactions", auth0.EnsureValidTokenMiddleware(auth0.AuthroizeMddlw(pkgHttp.GetTransactions)))
	router.POST("/v1/organizations", auth0.EnsureValidTokenMiddleware(auth0.AuthroizeMddlw(pkgHttp.CreateOrganization)))

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
