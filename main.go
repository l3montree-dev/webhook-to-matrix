package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/l3montree-dev/webhook-to-matrix/pkg/api"
)

const Port = 5001

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/webhook/glitchtip", api.GlitchTipHandler)
	addr := fmt.Sprintf("0.0.0.0:%d", Port)
	log.Printf("Listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
