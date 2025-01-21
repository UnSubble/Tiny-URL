package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/unsubble/tiny-url/common"
	"github.com/unsubble/tiny-url/config"
	"github.com/unsubble/tiny-url/database"
	"github.com/unsubble/tiny-url/tinyurl"
)

const CONFIG_PATH = "./config.json"

func newRepo(driverCfg *common.DriverConfig) database.Driver {
	switch driverCfg.DriverName {
	case "sqlite", "sqlite3":
		driver, err := database.NewSqliteRepository(driverCfg)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		return driver
	case "postgre", "postgres", "postgresql":
		driver, err := database.NewPostgresRepository(driverCfg)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		return driver
	default:
		log.Fatalf("Unknown database")
	}
	return nil
}

func shortenURLHandler(w http.ResponseWriter, r *http.Request, gen *tinyurl.Generator) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil || reqData.URL == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	tinyUrl, err := gen.GenerateURL(reqData.URL)
	if err != nil {
		log.Printf("[Error: %v", err)
	}

	respData := map[string]string{"short_url": tinyUrl.GeneratedURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respData)
}

func fetchOriginalURLHandler(w http.ResponseWriter, r *http.Request, driver database.Driver) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	shortURL := r.URL.Query().Get("short_url")
	if shortURL == "" {
		http.Error(w, "Missing 'short_url' parameter", http.StatusBadRequest)
		return
	}

	originalURL, err := driver.GetOriginalUrl(shortURL)

	if err != nil {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	respData := map[string]string{"original_url": originalURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respData)
}

func serveStaticPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func main() {
	cfg := config.ReadConfig(CONFIG_PATH)
	ctx := context.Background()
	repo := newRepo(cfg.DriverConfig)
	gen := tinyurl.NewGenerator(&ctx, cfg, repo, 0, 0)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", serveStaticPage)
	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		shortenURLHandler(w, r, gen)
	})
	http.HandleFunc("/fetch", func(w http.ResponseWriter, r *http.Request) {
		fetchOriginalURLHandler(w, r, repo)
	})

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
