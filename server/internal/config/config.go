package config

import (
    "time"
)

type Config struct {
    Port            string
    QuotesFile      string
    PowDifficulty   int
    ReadTimeout     time.Duration
    WriteTimeout    time.Duration
    MaxConnections  int
    ChallengeLength int
    ShutdownTimeout time.Duration
    // Новые параметры
    MaxRequestsPerIP    int           // Ограничение запросов с одного IP
    RateLimitWindow    time.Duration  // Окно для rate limiting
    BlacklistThreshold int           // Порог для занесения в черный список
    BlacklistDuration  time.Duration  // Время блокировки IP
    MaxMessageSize     int           // Максимальный размер сообщения
    SolutionTTL        time.Duration // Время жизни решения PoW
    BufferSize         int           // Размер буфера для чтения
}

func NewConfig() *Config {
    return &Config{
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
        MaxMessageSize:    1024,
        SolutionTTL:       5 * time.Minute,
        BufferSize:        1024,
    }
} 