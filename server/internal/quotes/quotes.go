package quotes

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"
)

// QuoteService предоставляет сервис для работы с цитатами
type QuoteService struct {
	// quotes содержит список всех доступных цитат
	quotes []string
}

// NewQuoteService создает новый экземпляр сервиса цитат
// filename - путь к JSON файлу с цитатами
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

// GetRandomQuote возвращает случайную цитату из списка
func (qs *QuoteService) GetRandomQuote() string {
	rand.Seed(time.Now().UnixNano())
	return qs.quotes[rand.Intn(len(qs.quotes))]
} 