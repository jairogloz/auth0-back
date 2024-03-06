package auth0

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Scope string `json:"scope"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c *CustomClaims) Validate(ctx context.Context) error {

	//token := ctx.Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	//
	fmt.Println(c)

	return nil

}

// EnsureValidToken is a middleware that will check the validity of our JWT.
func EnsureValidToken() func(next http.Handler) http.Handler {
	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH0_AUDIENCE")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Encountered error while validating JWT: %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Failed to validate JWT."}`))
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(next http.Handler) http.Handler {
		return middleware.CheckJWT(next)
	}
}

// HasScope checks whether our claims have a specific scope.
func (c *CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}

func EnsureValidTokenMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Convert the httprouter.Handle to http.Handler
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next(w, r, ps)
		})

		// Wrap the handler with the EnsureValidToken middleware
		validTokenHandler := EnsureValidToken()(nextHandler)

		// Serve the request with the new handler
		validTokenHandler.ServeHTTP(w, r)
	}
}

func AuthroizeMddlw(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

		reqPermissions := buildPermission(r.URL.String(), r.Method)
		fmt.Println(reqPermissions)

		claims := token.CustomClaims.(*CustomClaims)
		if !claims.HasScope(reqPermissions) {
			errResp := map[string]string{
				"message": fmt.Sprintf("Insufficient scope. Required: %s", reqPermissions),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			err := json.NewEncoder(w).Encode(errResp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		next(w, r, ps)
	}
}

// buildPermission returns a permission string based on the entity and action.
func buildPermission(rawURL string, httpVerb string) string {
	entityPart, err := getEntityPartFromURL(rawURL)
	if err != nil {
		return ""
	}

	action := mapHTTPVerbToAction(httpVerb)

	return fmt.Sprintf("%s:%s", action, entityPart)
}

func getEntityPartFromURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	pathSegments := strings.Split(path.Clean(parsedURL.Path), "/")
	for i, segment := range pathSegments {
		if segment == "v1" && i+1 < len(pathSegments) {
			return pathSegments[i+1], nil
		}
	}

	return "", nil
}

func mapHTTPVerbToAction(httpVerb string) string {
	switch httpVerb {
	case http.MethodGet:
		return "read"
	case http.MethodPost:
		return "create"
	case http.MethodPut:
		return "update"
	case http.MethodDelete:
		return "delete"
	default:
		return "unknown"
	}
}
