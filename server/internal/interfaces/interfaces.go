package interfaces

import (
	"context"
	"net"
	"time"
)

// QuoteProvider определяет интерфейс для сервиса цитат
type QuoteProvider interface {
	GetRandomQuote() string
}

// POWService определяет интерфейс для сервиса Proof of Work
type POWService interface {
	GenerateChallenge(length int) (string, error)
	VerifySolution(challenge, solution string) bool
}

// RateLimiter определяет интерфейс для контроля доступа по IP
type RateLimiter interface {
	IsAllowed(ip string) bool
}

// MetricsCollector определяет интерфейс для сбора метрик
type MetricsCollector interface {
	IncActiveConnections()
	DecActiveConnections()
	IncTotalConnections()
	IncFailedChallenges()
	IncSuccessChallenges()
	IncTotalQuotesSent()
	GetActiveConnections() int64
	GetStats() map[string]int64
	Wait()
}

// Logger определяет интерфейс для логирования
type Logger interface {
	Error(format string, v ...interface{})
	Info(format string, v ...interface{})
}

// Connection определяет интерфейс для сетевого соединения
type Connection interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
	RemoteAddr() net.Addr
	SetDeadline(t time.Time) error
}

// Server определяет интерфейс для сервера
type Server interface {
	Start(ctx context.Context) error
	Stop() error
}