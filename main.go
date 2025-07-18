package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/l3montree-dev/webhook-to-matrix/pkg/api"
)

const Port = 5001

func healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	// answer a json response with a 200 OK status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Check if required environment variables are set
	requiredVars := []string{"ROOM_ID", "ACCESS_TOKEN", "HOME_SERVER", "WEBHOOK_SECRET"}
	for _, v := range requiredVars {
		if value := os.Getenv(v); value == "" {
			log.Fatalf("Environment variable %s is not set", v)
		}
	}
	log.Println("All required environment variables are set.")

	http.HandleFunc("/health", healthCheckHandler)
	// set the webhook secret to all paths
	http.HandleFunc(fmt.Sprintf("/webhook/%s/glitchtip", os.Getenv("WEBHOOK_SECRET")), api.TransformGlitchTip)
	http.HandleFunc(fmt.Sprintf("/webhook/%s/botkube", os.Getenv("WEBHOOK_SECRET")), api.TransformBotKube)
	addr := fmt.Sprintf("0.0.0.0:%d", Port)
	log.Printf("Listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
