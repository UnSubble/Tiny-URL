package main

import (
	"fmt"

	config "github.com/unsubble/tiny-url/config"
)

const CONFIG_PATH = "./config.json"

func main() {

	cfg := config.ReadConfig(CONFIG_PATH)

	fmt.Println(cfg)

	// ctx := context.Background()

}
