package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func (app *Config) Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) < 1 {
			http.Error(w, "Not authorized", http.StatusInternalServerError)
			return
		}
		accessToken := strings.Split(authHeader, " ")[1]
		rptResult, err := app.GoCloakClient.RetrospectToken(r.Context(),
			accessToken, clientId, clientSecret, realm)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		isTokenValid := *rptResult.Active
		if !isTokenValid {
			http.Error(w, "Not authorized", http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *Config) ProtectForRole(next http.Handler, role string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) < 1 {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}
		accessToken := strings.Split(authHeader, " ")[1]
		rptResult, err := app.GoCloakClient.RetrospectToken(r.Context(),
			accessToken, clientId, clientSecret, realm)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		isTokenValid := *rptResult.Active
		if !isTokenValid {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		if !app.accessTokenContainsRole(r.Context(), accessToken, role) {
			http.Error(w, "Insuficient Previleges", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *Config) accessTokenContainsRole(ctx context.Context, accessToken string, roleName string) bool {
	_, claims, err := app.GoCloakClient.DecodeAccessToken(ctx, accessToken, realm)
	if err != nil {
		log.Println("Failed to decode the accessToken")
		return false
	}

	claimsMap := jwtClaimsToMap(claims)
	rolesMap, err := extractRolesFromMapClaims(claimsMap)

	if err != nil {
		log.Println("Failed to extract user permissions")
		return false
	}
	if !containsRole(rolesMap, roleName) {
		log.Println("Insuficient privileges")
		return false
	}
	return true
}

func jwtClaimsToMap(claims *jwt.MapClaims) map[string]any {
	customClaimMap := make(map[string]any)
	for key, val := range *claims {
		customClaimMap[key] = val
	}
	return customClaimMap
}

// / Creates a map in which the keys are the role names
// / To make it easier to do a lookup (O1) instead of loop through an array looking for the key
func extractRolesFromMapClaims(mapClaims map[string]any) (map[string]string, error) {
	roles, ok := mapClaims["realm_access"].(map[string]any)["roles"]
	if !ok {
		return nil, errors.New("Failed to get roles")
	}
	rolesMap := make(map[string]string)
	for _, v := range roles.([]any) {
		rolesMap[v.(string)] = ""
	}
	return rolesMap, nil
}

func containsRole(rolesMap map[string]string, role string) bool {
	_, ok := rolesMap[role]
	if !ok {
		return false
	}
	return true
}
