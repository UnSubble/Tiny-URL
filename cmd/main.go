package main

import (
	"context"
	"fmt"
	"log"

	"github.com/unsubble/tiny-url/config"
	"github.com/unsubble/tiny-url/database"
	"github.com/unsubble/tiny-url/tinyurl"
)

const CONFIG_PATH = "./config.json"

func main() {
	ctx := context.Background()

	cfg := config.ReadConfig(CONFIG_PATH)

	repo, err := database.NewPostgresRepository(cfg.DriverConfig[1])

	if err != nil {
		log.Fatalf("%v", err)
	}

	generator := tinyurl.NewGenerator(&ctx, cfg, repo, 0, 0)
	url, err := generator.GenerateURL("example.com/")

	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("url:", url)
}
