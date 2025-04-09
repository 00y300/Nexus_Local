package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

// PageData holds the data passed to the HTML template.
type PageData struct {
	LoginURL    string
	RedirectURL string
}

// App holds the OAuth configuration and the HTML template.
type App struct {
	OAuthCfg *oauth2.Config
	Tmpl     *template.Template
}

// NewApp creates a new App instance for handling OAuth.
func NewApp(oauthCfg *oauth2.Config, tmpl *template.Template) *App {
	return &App{
		OAuthCfg: oauthCfg,
		Tmpl:     tmpl,
	}
}

// Root handler: renders the index page with a link to log in.
func (a *App) Root(w http.ResponseWriter, r *http.Request) {
	state := "some-random-state"
	loginURL := a.OAuthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)

	data := PageData{
		LoginURL:    loginURL,
		RedirectURL: a.OAuthCfg.RedirectURL,
	}

	if err := a.Tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Login handler: redirects the user to Microsoft sign-in.
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	state := "some-random-state"
	url := a.OAuthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	log.Printf("Redirecting user to: %s", url)
	http.Redirect(w, r, url, http.StatusFound)
}

// OAuthCallback handler: processes the authorization code from Microsoft.
func (a *App) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code parameter", http.StatusBadRequest)
		return
	}
	log.Printf("Received code: %s", code)

	token, err := a.OAuthCfg.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if token.RefreshToken != "" {
		log.Printf("Refresh Token: %s", token.RefreshToken)
	} else {
		log.Printf("No Refresh Token returned.")
	}
	log.Printf("Access Token: %s", token.AccessToken)

	if idToken, ok := token.Extra("id_token").(string); ok {
		log.Printf("ID Token: %s", idToken)
	} else {
		log.Printf("No ID Token returned.")
	}

	// Call Microsoft Graph /me endpoint
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)
	if err != nil {
		http.Error(w, "Failed to create Graph request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get Graph response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Unexpected status code from Graph: %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	var userData map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		http.Error(w, "Failed to decode JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	prettyJSON, err := json.MarshalIndent(userData, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(prettyJSON)
}
