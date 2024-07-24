package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Quote struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request to server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Server response error: %s", resp.Status)
	}

	var quote Quote
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	err = os.WriteFile("quote.txt", []byte("Dólar: "+quote.Bid), 0644)
	if err != nil {
		log.Fatalf("Error saving quote to file: %v", err)
	}

	log.Println("Quote saved successfully to quote.txt")
}
