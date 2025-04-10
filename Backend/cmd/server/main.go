// cmd/server/main.go
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
	// ensure uploads dir exists
	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Fatalf("could not create uploads dir: %v", err)
	}

	// 1) Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on real ENV")
	}

	// 2) Azure AD settings...
	tenantID := os.Getenv("AZUREAD_TENANT_ID")
	clientID := os.Getenv("AZUREAD_APP_ID")
	clientSecret := os.Getenv("AZUREAD_VALUE")
	if tenantID == "" || clientID == "" || clientSecret == "" {
		log.Fatal("Missing one of AZUREAD_TENANT_ID, AZUREAD_APP_ID, AZUREAD_VALUE")
	}

	// 3) OAuth2 config
	oauthCfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/redirect",
		Scopes:       []string{"openid", "profile", "email", "offline_access", "User.Read.All"},
	}

	// 4) OIDC provider & verifier
	ctx := context.Background()
	issuer := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", tenantID)
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Fatalf("Failed to initialize OIDC provider: %v", err)
	}
	oauthCfg.Endpoint = provider.Endpoint()
	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})

	// 5) Template
	tmpl := template.Must(
		template.New("index.html").
			Funcs(template.FuncMap{"Join": strings.Join}).
			ParseFiles("templates/index.html"),
	)

	// 6) AuthApp
	authApp := auth.NewApp(oauthCfg, verifier, tmpl)

	// 7) DB
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

	// 8) Server
	srv := server.NewServer(authApp, sqlDB)
	log.Println("🚀 Starting server on :8080")
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
