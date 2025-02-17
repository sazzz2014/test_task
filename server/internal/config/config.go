package config

import (
    "fmt"
    "time"
)

type Config struct {
    // Базовые настройки сервера
    Port            string          // Порт для прослушивания
    QuotesFile      string          // Путь к файлу с цитатами
    PowDifficulty   int             // Сложность Proof of Work
    
    // Таймауты
    ReadTimeout     time.Duration   // Таймаут на чтение
    WriteTimeout    time.Duration   // Таймаут на запись
    ShutdownTimeout time.Duration   // Таймаут на graceful shutdown
    
    // Ограничения соединений
    MaxConnections  int             // Максимальное число одновременных соединений
    ChallengeLength int             // Длина challenge для PoW
    
    // Rate limiting и защита
    MaxRequestsPerIP    int         // Максимум запросов с одного IP
    RateLimitWindow    time.Duration // Временное окно для rate limiting
    BlacklistThreshold int          // После скольких нарушений IP попадает в черный список
    BlacklistDuration  time.Duration // На какое время IP блокируется
    
    // Ограничения буферов
    MaxMessageSize     int          // Максимальный размер сообщения
    BufferSize         int          // Размер буфера для чтения
    
    // PoW настройки
    SolutionTTL        time.Duration // Время жизни решения PoW
}
// Проверяет валидность конфигурации
func (c *Config) Validate() error {
    if c.Port == "" {
        return fmt.Errorf("port must be specified")
    }
    if c.MaxConnections <= 0 {
        return fmt.Errorf("max connections must be positive")
    }
    if c.ReadTimeout <= 0 || c.WriteTimeout <= 0 {
        return fmt.Errorf("timeouts must be positive")
    }
    if c.BufferSize <= 0 || c.MaxMessageSize <= 0 {
        return fmt.Errorf("buffer sizes must be positive")
    }
    return nil
}

func NewConfig() *Config {
    cfg := &Config{
        Port:              ":8080",
        QuotesFile:        "quotes.json",
        PowDifficulty:     4,
        ReadTimeout:       30 * time.Second,
        WriteTimeout:      30 * time.Second,
        MaxConnections:    100,
        ChallengeLength:   32,
        ShutdownTimeout:   30 * time.Second,
        MaxRequestsPerIP:  100,
        RateLimitWindow:   time.Minute,
        BlacklistThreshold: 5,
        BlacklistDuration: 24 * time.Hour,
        MaxMessageSize:    1024,    // Ограничиваем размер сообщения
        SolutionTTL:       5 * time.Minute,
        BufferSize:        1024,    // Ограничиваем размер буфера
    }
    if err := cfg.Validate(); err != nil {
        panic(fmt.Sprintf("invalid config: %v", err))
    }
    return cfg
} 