package auth

import (
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/coreos/go-oidc"
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
	Verifier *oidc.IDTokenVerifier
	Tmpl     *template.Template
}

// NewApp creates a new App.
//
//	func NewApp(oauthCfg *oauth2.Config, tmpl *template.Template) *App {
//		return &App{OAuthCfg: oauthCfg, Tmpl: tmpl}
//	}
// func NewApp(ctx context.Context, oauthCfg *oauth2.Config, issuer, clientID string) (*App, error) {
// 	provider, err := oidc.NewProvider(ctx, issuer)
// 	if err != nil {
// 		return nil, err
// 	}
// 	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
// 	return &App{OAuthCfg: oauthCfg, Verifier: verifier}, nil
// }

// in internal/auth/auth.go
func NewApp(oauthCfg *oauth2.Config, verifier *oidc.IDTokenVerifier, tmpl *template.Template) *App {
	return &App{
		OAuthCfg: oauthCfg,
		Verifier: verifier,
		Tmpl:     tmpl,
	}
}

// auth.go (continued)

var (
	ErrNoAuthHeader = errors.New("no Authorization header")
	ErrInvalidToken = errors.New("invalid token")
	ErrForbidden    = errors.New("forbidden")
)

// AuthMiddleware ensures the request has a valid Bearer token
// and that the “roles” claim includes “admin”.
func (a *App) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hdr := r.Header.Get("Authorization")
		if hdr == "" {
			http.Error(w, ErrNoAuthHeader.Error(), http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(hdr, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}
		rawIDToken := parts[1]

		// Verify signature & claims (exp, aud, iss, ...)
		idToken, err := a.Verifier.Verify(r.Context(), rawIDToken)
		if err != nil {
			http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}

		// Pull out the roles claim
		var claims struct {
			Roles []string `json:"roles"`
		}
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "failed to parse claims", http.StatusInternalServerError)
			return
		}

		// Check for “admin”
		isAdmin := slices.Contains(claims.Roles, "admin")
		if !isAdmin {
			http.Error(w, ErrForbidden.Error(), http.StatusForbidden)
			return
		}

		// OK to proceed
		next.ServeHTTP(w, r)
	})
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

	// Extract the ID token string
	idt, _ := token.Extra("id_token").(string)
	if idt == "" {
		http.Error(w, "No id_token in response", http.StatusInternalServerError)
		return
	}

	// Set it as a secure, HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "id_token",
		Value:    idt,
		Path:     "/",
		HttpOnly: true,
		// Secure: true, // in prod
		// SameSite: http.SameSiteLaxMode,
	})

	// Redirect back to your Next.js admin page (or home)
	http.Redirect(w, r, "http://localhost:3000/admin/add-item", http.StatusFound)
}
