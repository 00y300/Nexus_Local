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

type PageData struct {
	LoginURL    string
	RedirectURL string
}

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

// OAuthCallback handler: processes the authorization code from Microsoft,
// exchanges it for an access token, and then calls Microsoft Graph /me.
func (a *App) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code parameter", http.StatusBadRequest)
		return
	}
	log.Printf("Received code: %s", code)

	// Exchange the authorization code for a token.
	token, err := a.OAuthCfg.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Log refresh and access tokens (for debugging/demo).
	if token.RefreshToken != "" {
		log.Printf("Refresh Token: %s", token.RefreshToken)
	} else {
		log.Printf("No Refresh Token returned.")
	}
	log.Printf("Access Token: %s", token.AccessToken)

	// Check for an ID token (OpenID Connect flow).
	if idToken, ok := token.Extra("id_token").(string); ok {
		log.Printf("ID Token: %s", idToken)
	} else {
		log.Printf("No ID Token returned.")
	}

	//-----------------------------------------------------------------
	// 1. Create a new request to the Microsoft Graph /me endpoint.
	//-----------------------------------------------------------------
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)
	if err != nil {
		http.Error(w, "Failed to create Graph request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Add the access token to the request header.
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// 3. Send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get Graph response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 4. Check the response status.
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Unexpected status code from Graph: %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	// 5. Decode the JSON response into a map (or custom struct).
	var userData map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		http.Error(w, "Failed to decode JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 6. Pretty‐print the JSON response and send it back to the user.
	prettyJSON, err := json.MarshalIndent(userData, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(prettyJSON)
}
