package main

import (
	"fmt"
	"log"
	"time"

	"client/internal/client"
)

func main() {
	c := client.NewClient(":8080", 4, 30*time.Second)
	
	quote, err := c.GetQuote()
	if err != nil {
		log.Fatalf("Error getting quote: %v", err)
	}

	fmt.Printf("Received quote: %s\n", quote)
} 