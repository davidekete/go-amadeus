package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	origins := []string{"PAR", "ORY", "MAD"}
	results := FetchDestinationsAsync(origins, "300")

	for msg := range results {
		switch v := msg.(type) {
		case *FlightResponse:
			pretty, _ := json.MarshalIndent(v, "", "  ")
			fmt.Println(string(pretty))
		case error:
			log.Printf("Request failed: %v\n", v)
		}
	}
}
