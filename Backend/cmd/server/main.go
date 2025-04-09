package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"

	"nexus.local/internal/auth"
	"nexus.local/internal/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, proceeding without it.")
	}

	tenantID := os.Getenv("AZUREAD_TENANT_ID")
	if tenantID == "" {
		log.Fatal("Missing AZUREAD_TENANT_ID in environment.")
	}
	clientID := os.Getenv("AZUREAD_APP_ID")
	if clientID == "" {
		log.Fatal("Missing AZUREAD_APP_ID in environment.")
	}
	clientSecret := os.Getenv("AZUREAD_VALUE")
	if clientSecret == "" {
		log.Fatal("Missing AZUREAD_VALUE in environment.")
	}

	scopes := []string{"offline_access", "User.Read.All"}
	log.Printf("Using scopes: %v", scopes)

	endpoint := oauth2.Endpoint{
		AuthURL:  fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", tenantID),
		TokenURL: fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID),
	}

	oauthCfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     endpoint,
		RedirectURL:  "http://localhost:8080/redirect",
		Scopes:       scopes,
	}

	tmpl := template.Must(
		template.New("index.html").Funcs(template.FuncMap{"Join": strings.Join}).ParseFiles("templates/index.html"),
	)

	authApp := auth.NewApp(oauthCfg, tmpl)
	srv := server.NewServer(authApp)

	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
