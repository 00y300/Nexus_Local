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

// App holds the OAuth config and template.
type App struct {
	OAuthCfg *oauth2.Config
	Tmpl     *template.Template
}

// NewApp creates a new App.
func NewApp(oauthCfg *oauth2.Config, tmpl *template.Template) *App {
	return &App{OAuthCfg: oauthCfg, Tmpl: tmpl}
}

// Root renders the login page.
func (a *App) Root(w http.ResponseWriter, r *http.Request) {
	state := "some-random-state"
	loginURL := a.OAuthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	data := PageData{LoginURL: loginURL, RedirectURL: a.OAuthCfg.RedirectURL}
	if err := a.Tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Login redirects to Azure AD.
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	state := "some-random-state"
	url := a.OAuthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	log.Printf("Redirecting to: %s", url)
	http.Redirect(w, r, url, http.StatusFound)
}

// OAuthCallback handles the code exchange and calls Graph /me.
func (a *App) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code", http.StatusBadRequest)
		return
	}
	token, err := a.OAuthCfg.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Token exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Access Token: %s", token.AccessToken)
	if rt := token.RefreshToken; rt != "" {
		log.Printf("Refresh Token: %s", rt)
	}
	if idt, ok := token.Extra("id_token").(string); ok {
		log.Printf("ID Token: %s", idt)
	}

	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)
	if err != nil {
		http.Error(w, "Graph request creation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Graph request failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Graph returned %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}
	var user map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		http.Error(w, "JSON decode failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	pretty, _ := json.MarshalIndent(user, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(pretty)
}
