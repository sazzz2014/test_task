package quotes

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"
)

type QuoteService struct {
	quotes []string
}

func NewQuoteService(filename string) (*QuoteService, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var quotes []string
	if err := json.Unmarshal(data, &quotes); err != nil {
		return nil, err
	}

	return &QuoteService{quotes: quotes}, nil
}

func (qs *QuoteService) GetRandomQuote() string {
	rand.Seed(time.Now().UnixNano())
	return qs.quotes[rand.Intn(len(qs.quotes))]
} 