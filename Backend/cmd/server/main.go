package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/coreos/go-oidc"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"

	"nexus.local/internal/auth"
	"nexus.local/internal/db"
	"nexus.local/internal/server"
)

func main() {
	// 1) Load .env (if present)
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on real ENV")
	}

	// 2) Pull Azure AD settings
	tenantID := os.Getenv("AZUREAD_TENANT_ID")
	clientID := os.Getenv("AZUREAD_APP_ID")
	clientSecret := os.Getenv("AZUREAD_VALUE")
	if tenantID == "" || clientID == "" || clientSecret == "" {
		log.Fatal("Missing one of AZUREAD_TENANT_ID, AZUREAD_APP_ID, AZUREAD_VALUE")
	}

	// 3) Build OAuth2 config (we’ll fill Endpoint from the OIDC provider)
	oauthCfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/redirect",
		Scopes: []string{
			"openid", // to get an ID token
			"profile",
			"email",
			"offline_access",
			"User.Read.All", // your Graph scope
		},
	}

	// 4) Set up OIDC provider & verifier
	ctx := context.Background()
	issuer := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", tenantID)
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Fatalf("Failed to initialize OIDC provider: %v", err)
	}
	// fill in the AuthURL/TokenURL automatically
	oauthCfg.Endpoint = provider.Endpoint()

	// verifier will check signature, expiry, audience, issuer, etc.
	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})

	// 5) Parse your index.html template (for the “/” login page)
	tmpl := template.Must(
		template.New("index.html").
			Funcs(template.FuncMap{"Join": strings.Join}).
			ParseFiles("templates/index.html"),
	)

	// 6) Initialize your auth.App (with OAuth2 + OIDC + template)
	authApp := auth.NewApp(oauthCfg, verifier, tmpl)

	// 7) Connect to your database
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	sqlDB, err := db.Connect(dbUser, dbPass, dbHost, dbPort, dbName)
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}
	defer sqlDB.Close()
	log.Println("✅ Connected to database.")

	// 8) Wire up the HTTP server
	srv := server.NewServer(authApp, sqlDB)
	log.Println("🚀 Starting server on :8080")
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
