package server

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"server/internal/config"
	"server/internal/interfaces"
	"server/internal/protocol"
)

type Server struct {
	config       *config.Config
	pow          interfaces.POWService
	quoteService interfaces.QuoteProvider
	listener     net.Listener
	connections  sync.Map
	shutdown     chan struct{}
	metrics      interfaces.MetricsCollector
	ipControl    interfaces.RateLimiter
	logger       interfaces.Logger
}

func NewServer(
	cfg *config.Config,
	pow interfaces.POWService,
	qs interfaces.QuoteProvider,
	metrics interfaces.MetricsCollector,
	ipControl interfaces.RateLimiter,
	logger interfaces.Logger,
) *Server {
	return &Server{
		config:       cfg,
		pow:          pow,
		quoteService: qs,
		shutdown:     make(chan struct{}),
		logger:       logger,
		metrics:      metrics,
		ipControl:    ipControl,
	}
}

func (s *Server) Start(ctx context.Context) error {
	var err error
	s.listener, err = net.Listen("tcp", s.config.Port)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	s.logger.Info("Server started on %s", s.config.Port)

	go func() {
		<-ctx.Done()
		s.logger.Info("Initiating graceful shutdown...")
		s.Stop()
	}()

	go s.acceptConnections(ctx)

	<-ctx.Done()
	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("Stopping server...")
	
	select {
	case <-s.shutdown: // Проверяем, не закрыт ли уже канал
		// Канал уже закрыт, ничего не делаем
	default:
		close(s.shutdown)
	}

	// Закрываем listener для прекращения приема новых соединений
	if s.listener != nil {
		s.listener.Close()
	}

	// Ожидаем завершения всех активных соединений
	done := make(chan struct{})
	go func() {
		s.metrics.Wait() // Ждем пока все соединения завершатся
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("All connections closed gracefully")
	case <-time.After(s.config.ShutdownTimeout):
		s.logger.Info("Shutdown timeout exceeded, forcing shutdown")
	}

	return nil
}

func (s *Server) acceptConnections(ctx context.Context) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.shutdown:
				return
			default:
				if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
					time.Sleep(time.Millisecond * 100)
					continue
				}
				s.logger.Error("Accept error: %v", err)
				return
			}
		}

		// Проверяем лимит подключений
		currentConnections := s.metrics.GetActiveConnections()
		if currentConnections >= int64(s.config.MaxConnections) {
			s.logger.Info("Connection limit reached, rejecting connection from %s", conn.RemoteAddr())
			conn.Close()
			continue
		}

		go s.handleConnection(ctx, conn)
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		conn.Close()
		s.metrics.DecActiveConnections()
	}()

	// Получаем IP клиента
	ip := conn.RemoteAddr().(*net.TCPAddr).IP.String()

	// Проверяем ограничения по IP
	if !s.ipControl.IsAllowed(ip) {
		s.logger.Info("Connection rejected from %s (rate limit/blacklist)", ip)
		return
	}

	s.metrics.IncTotalConnections()
	s.metrics.IncActiveConnections()

	// Устанавливаем таймауты
	deadline := time.Now().Add(s.config.ReadTimeout)
	conn.SetDeadline(deadline)

	s.connections.Store(conn.RemoteAddr(), conn)
	defer s.connections.Delete(conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	// Ожидаем HELLO
	message, err := reader.ReadString('\n')
	if err != nil || !strings.HasPrefix(message, protocol.CmdHello) {
		return
	}

	// Генерируем и отправляем вызов
	challenge, err := s.pow.GenerateChallenge(s.config.ChallengeLength)
	if err != nil {
		s.logger.Error("Failed to generate challenge: %v", err)
		return
	}
	fmt.Fprintf(conn, "%s %s\n", protocol.CmdChallenge, challenge)

	// Ожидаем решение
	message, err = reader.ReadString('\n')
	if err != nil {
		return
	}

	parts := strings.Fields(message)
	if len(parts) != 2 || parts[0] != protocol.CmdSolution {
		return
	}

	// Проверяем решение
	if s.pow.VerifySolution(challenge, parts[1]) {
		s.metrics.IncSuccessChallenges()
		quote := s.quoteService.GetRandomQuote()
		s.metrics.IncTotalQuotesSent()
		fmt.Fprintf(conn, "%s %s\n", protocol.CmdQuote, quote)
	} else {
		s.metrics.IncFailedChallenges()
		fmt.Fprintf(conn, "%s\n", protocol.CmdError)
	}
}
