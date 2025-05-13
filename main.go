package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	response, err := GetRequest("PAR", "300")

	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(response, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
