package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"nexus.local/cmd/server"
	"nexus.local/internal/auth"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	// Load environment variables (if .env is present).
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, proceeding without it.")
	}

	// Read required environment variables.
	tenantID := os.Getenv("AZUREAD_TENANT_ID")
	if tenantID == "" {
		log.Fatal("Missing AZUREAD_TENANT_ID in environment.")
	}
	microsoftAppID := os.Getenv("AZUREAD_APP_ID")
	if microsoftAppID == "" {
		log.Fatal("Missing AZUREAD_APP_ID in environment.")
	}
	microsoftSecretValue := os.Getenv("AZUREAD_VALUE")
	if microsoftSecretValue == "" {
		log.Fatal("Missing AZUREAD_VALUE in environment.")
	}

	// Define the scopes for your Azure app.
	allScopes := []string{"offline_access", "User.Read.All"}
	log.Printf("Using scopes: %v", allScopes)

	// Construct the OAuth 2.0 endpoints using the tenant ID.
	endpoint := oauth2.Endpoint{
		AuthURL:  fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", tenantID),
		TokenURL: fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID),
	}

	// Create the OAuth2 config (must match port 8080 for the redirect).
	oauthCfg := &oauth2.Config{
		ClientID:     microsoftAppID,
		ClientSecret: microsoftSecretValue,
		Endpoint:     endpoint,
		RedirectURL:  "http://localhost:8080/redirect",
		Scopes:       allScopes,
	}

	// Parse our HTML template and provide a Join function (if needed).
	tmpl := template.Must(
		template.New("index.html").
			Funcs(template.FuncMap{"Join": strings.Join}).
			ParseFiles("index.html"),
	)

	// Create our OAuth "App" with its configuration and template.
	oauthApp := auth.NewApp(oauthCfg, tmpl)

	// Create a new server instance using dependency injection.
	srv := server.NewServer(oauthApp)

	// Start the server on port 8080.
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
