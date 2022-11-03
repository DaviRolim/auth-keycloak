package main

import (
	"log"
	"net/http"

	"github.com/Nerzal/gocloak/v12"
)

var (
	clientId     = "gosimpleapi"
	KeycloakURL  = "http://localhost:8080"
	clientSecret = "xArtw6iv0CWAdYDxAA1PABBAhjFKZ4zn"
	realm        = "goapi"
)

type Config struct {
	GoCloakClient *gocloak.GoCloak
}

func main() {
	app := Config{
		GoCloakClient: gocloak.NewClient(KeycloakURL),
	}

	log.Fatal(http.ListenAndServe(":8081", app.routes()))

}
