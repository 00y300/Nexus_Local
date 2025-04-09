package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"

	"nexus.local/internal/auth"
	"nexus.local/internal/db"
	"nexus.local/internal/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found.")
	}

	// --- OAuth config ---
	tenantID := os.Getenv("AZUREAD_TENANT_ID")
	clientID := os.Getenv("AZUREAD_APP_ID")
	clientSecret := os.Getenv("AZUREAD_VALUE")
	if tenantID == "" || clientID == "" || clientSecret == "" {
		log.Fatal("Missing one of AZUREAD_TENANT_ID, AZUREAD_APP_ID, AZUREAD_VALUE")
	}
	scopes := []string{"offline_access", "User.Read.All"}
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
		template.New("index.html").
			Funcs(template.FuncMap{"Join": strings.Join}).
			ParseFiles("templates/index.html"),
	)

	// --- DB config & connect ---
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
	log.Println("Connected to database.")

	// (Optional) list items at startup
	items, _ := db.GetAllItems(sqlDB)
	for _, it := range items {
		log.Printf("Item %d: %s ($%.2f) stock=%d",
			it.ID, it.Name, it.Price, it.Stock)
	}

	// --- Wire up server ---
	authApp := auth.NewApp(oauthCfg, tmpl)
	srv := server.NewServer(authApp, sqlDB)

	log.Println("Starting server on :8080")
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
