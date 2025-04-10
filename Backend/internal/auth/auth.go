package auth

import (
	"context"
	"database/sql"
	"errors"
	"html/template"
	"net/http"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

// ctxKey is the type we use for context keys in this package
type ctxKey string

// ContextKeyUser is the context key under which we store the user's OID
const ContextKeyUser ctxKey = "userID"

// PageData holds the data passed to the HTML template.
type PageData struct {
	LoginURL    string
	RedirectURL string
}

// App holds the OAuth2 config, OIDC verifier, HTML template—and now a DB handle.
type App struct {
	OAuthCfg *oauth2.Config
	Verifier *oidc.IDTokenVerifier
	Tmpl     *template.Template
	DB       *sql.DB
}

// NewApp constructs a new App.
func NewApp(oauthCfg *oauth2.Config, verifier *oidc.IDTokenVerifier, tmpl *template.Template, db *sql.DB) *App {
	return &App{
		OAuthCfg: oauthCfg,
		Verifier: verifier,
		Tmpl:     tmpl,
		DB:       db,
	}
}

var (
	ErrNoAuthHeader = errors.New("no authorization cookie")
	ErrInvalidToken = errors.New("invalid token")
	ErrForbidden    = errors.New("forbidden")
)

// AuthMiddleware verifies the cookie, looks up is_admin in MySQL, and rejects non‑admins.
func (a *App) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1) Grab the ID token from the cookie
		ck, err := r.Cookie("id_token")
		if err != nil {
			http.Error(w, ErrNoAuthHeader.Error(), http.StatusUnauthorized)
			return
		}

		// 2) Verify signature & standard claims
		rawID := ck.Value
		idToken, err := a.Verifier.Verify(r.Context(), rawID)
		if err != nil {
			http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}

		// 3) Extract the Azure AD user OID
		var claims struct {
			OID string `json:"oid"`
		}
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "failed to parse token claims", http.StatusInternalServerError)
			return
		}

		// 4) Look up is_admin in your users table
		var isAdmin bool
		err = a.DB.QueryRowContext(
			r.Context(),
			"SELECT is_admin FROM users WHERE id = ?",
			claims.OID,
		).Scan(&isAdmin)

		if err == sql.ErrNoRows {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !isAdmin {
			http.Error(w, ErrForbidden.Error(), http.StatusForbidden)
			return
		}

		// 5) Inject the user ID into context for downstream handlers
		ctx := context.WithValue(r.Context(), ContextKeyUser, claims.OID)
		next.ServeHTTP(w, r.WithContext(ctx))
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

// Login redirects the user to Azure AD for authentication.
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	state := "some-random-state"
	url := a.OAuthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

// OAuthCallback handles the code exchange, logs scopes, sets ID and access token cookies,
// and redirects back to your Next.js admin page.
func (a *App) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	token, err := a.OAuthCfg.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "token exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract raw ID token + access token
	idt, _ := token.Extra("id_token").(string)
	at := token.AccessToken

	// Set them as HttpOnly cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "id_token",
		Value:    idt,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    at,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	// Redirect back to your front‑end
	http.Redirect(w, r, "http://localhost:3000/admin/add-item", http.StatusFound)
}
