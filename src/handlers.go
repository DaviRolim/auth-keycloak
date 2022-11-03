package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type LoginFields struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (app *Config) Login(w http.ResponseWriter, r *http.Request) {
	var loginFields LoginFields
	ctx := context.Background()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error Getting Body", http.StatusInternalServerError)
	}
	json.Unmarshal(reqBody, &loginFields)
	log.Println(loginFields)
	jwt, err := app.GoCloakClient.Login(ctx, clientId, clientSecret, realm, loginFields.Username, loginFields.Password)
	if err != nil || jwt == nil {
		http.Error(w, "Error on login", http.StatusInternalServerError)
	}

	payload := map[string]any{
		"Error":       false,
		"Message":     fmt.Sprintf("User Authenticated Successfullly"),
		"AccessToken": jwt.AccessToken,
	}
	out, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(out)
}

func (app *Config) Greet(w http.ResponseWriter, r *http.Request) {
	payload := map[string]any{
		"Error":   false,
		"Message": "Greetings friend",
		"Data":    "",
	}
	out, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(out)
}

func (app *Config) GreetVip(w http.ResponseWriter, r *http.Request) {
	payload := map[string]any{
		"Error":   false,
		"Message": "Greetings friend VIP",
		"Data":    "",
	}
	out, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(out)
}
