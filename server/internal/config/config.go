package config

import (
    "fmt"
    "time"
)

// Config содержит все настройки конфигурации сервера
type Config struct {
    // Port определяет порт, на котором будет работать сервер
    Port            string          // Порт для прослушивания
    // QuotesFile путь к JSON файлу с цитатами
    QuotesFile      string          // Путь к файлу с цитатами
    // PowDifficulty задает сложность Proof of Work (количество начальных нулевых битов)
    PowDifficulty   int             // Сложность Proof of Work
    
    // ReadTimeout максимальное время ожидания чтения данных от клиента
    ReadTimeout     time.Duration   // Таймаут на чтение
    // WriteTimeout максимальное время ожидания записи данных клиенту
    WriteTimeout    time.Duration   // Таймаут на запись
    // ShutdownTimeout время ожидания завершения всех соединений при остановке сервера
    ShutdownTimeout time.Duration   // Таймаут на graceful shutdown
    
    // MaxConnections максимальное количество одновременных подключений к серверу
    MaxConnections  int             // Максимальное число одновременных соединений
    // ChallengeLength длина генерируемого challenge в байтах
    ChallengeLength int             // Длина challenge для PoW
    
    // MaxRequestsPerIP максимальное количество запросов с одного IP в течение RateLimitWindow
    MaxRequestsPerIP    int         // Максимум запросов с одного IP
    // RateLimitWindow временное окно для подсчета запросов с одного IP
    RateLimitWindow    time.Duration // Временное окно для rate limiting
    // BlacklistThreshold количество нарушений, после которых IP попадает в черный список
    BlacklistThreshold int          // После скольких нарушений IP попадает в черный список
    // BlacklistDuration время, на которое IP блокируется при попадании в черный список
    BlacklistDuration  time.Duration // На какое время IP блокируется
    
    // MaxMessageSize максимальный размер сообщения от клиента в байтах
    MaxMessageSize     int          // Максимальный размер сообщения
    // BufferSize размер буфера для чтения сообщений в байтах
    BufferSize         int          // Размер буфера для чтения
    
    // SolutionTTL время жизни решения PoW
    SolutionTTL        time.Duration // Время жизни решения PoW
}

// Validate проверяет корректность настроек конфигурации
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

// NewConfig создает новый экземпляр конфигурации с значениями по умолчанию
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