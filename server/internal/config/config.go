package config

type Config struct {
    Port            string
    QuotesFile      string
    PowDifficulty   int
    ReadTimeout     int
    WriteTimeout    int
    MaxConnections  int
    ChallengeLength int
}

func NewConfig() *Config {
    return &Config{
        Port:            ":8080",
        QuotesFile:      "quotes.json",
        PowDifficulty:   4,            // количество нулевых битов
        ReadTimeout:     30,           // секунды
        WriteTimeout:    30,           // секунды
        MaxConnections:  100,
        ChallengeLength: 32,
    }
} 