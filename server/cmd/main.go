package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"server/internal/config"
	"server/internal/logger"
	"server/internal/metrics"
	"server/internal/pow"
	"server/internal/quotes"
	"server/internal/ratelimit"
	"server/internal/server"
)

func main() {
	cfg := config.NewConfig()
	
	powService := pow.NewPoW(cfg.PowDifficulty)
	
	quoteService, err := quotes.NewQuoteService(cfg.QuotesFile)
	if err != nil {
		log.Fatalf("Failed to initialize quote service: %v", err)
	}

	// Создаем все необходимые сервисы
	loggerService := logger.NewLogger()
	metricsService := metrics.NewMetrics()
	ipControlService := ratelimit.NewIPControl(cfg)

	srv := server.NewServer(
		cfg,
		powService,
		quoteService,
		metricsService,
		ipControlService,
		loggerService,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработка сигналов для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
} 