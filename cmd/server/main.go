package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/hx23840/YunDoBridge/internal/config"
	handler "github.com/hx23840/YunDoBridge/internal/http"
	"github.com/hx23840/YunDoBridge/internal/types"
)

func main() {
	configPath := flag.String("config", "config.json", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	env := &types.Env{
		OpenAIAPIKey:        cfg.OpenAIAPIKey,
		OpenAIModelEndpoint: cfg.OpenAIModelEndpoint,
		CallsBaseURL:        cfg.CallsBaseURL,
		CallsAppID:          cfg.CallsAppID,
		CallsAppToken:       cfg.CallsAppToken,
	}

	http.HandleFunc("/endpoint", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleRequest(w, r, env)
	})

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
