package main

import (
	"log"

	"github.com/liyown/search4ai-go/api"
)

func main() {
	log.Println("Starting search4ai-go server...")
	if err := api.StartServer(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
