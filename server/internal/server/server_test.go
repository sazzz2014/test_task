package server

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"server/internal/config"
	"server/internal/mocks"
)

// fakeTCPConn оборачивает net.Conn для тестов
type fakeTCPConn struct {
	net.Conn
}

func (f fakeTCPConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345}
}

// Создаем вспомогательную функцию для настройки тестового окружения
func setupTest(t *testing.T) (*Server, *mocks.MockPOWService, *mocks.MockQuoteProvider, *mocks.MockMetricsCollector, *mocks.MockRateLimiter, *mocks.MockLogger) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPow := mocks.NewMockPOWService(ctrl)
	mockQuotes := mocks.NewMockQuoteProvider(ctrl)
	mockMetrics := mocks.NewMockMetricsCollector(ctrl)
	mockRateLimiter := mocks.NewMockRateLimiter(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)

	cfg := &config.Config{
		Port:            ":0",
		ReadTimeout:     time.Second,
		WriteTimeout:    time.Second,
		MaxConnections:  10,
		ChallengeLength: 32,
		ShutdownTimeout: time.Second,
		BufferSize:      1024,
		MaxMessageSize:  1024,
	}

	srv := NewServer(cfg, mockPow, mockQuotes, mockMetrics, mockRateLimiter, mockLogger)

	return srv, mockPow, mockQuotes, mockMetrics, mockRateLimiter, mockLogger
}

func TestServer_Start(t *testing.T) {
	srv, _, _, mockMetrics, _, mockLogger := setupTest(t)

	// Настраиваем ожидаемое поведение моков
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockMetrics.EXPECT().Wait().AnyTimes()

	// Запускаем сервер в отдельной горутине
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start(ctx)
	}()

	// Даем серверу время на запуск
	time.Sleep(100 * time.Millisecond)

	// Проверяем, что сервер успешно запустился
	require.NotNil(t, srv.listener)
	
	// Отменяем контекст, чтобы остановить сервер
	cancel()

	// Проверяем, что сервер корректно завершился
	err := <-errCh
	assert.NoError(t, err)
}

func TestServer_HandleConnection_Success(t *testing.T) {
	srv, mockPow, mockQuotes, mockMetrics, mockRateLimiter, mockLogger := setupTest(t)

	// Создаем пару тестовых соединений
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()
	server = fakeTCPConn{Conn: server}

	// Настраиваем ожидания для всех моков
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockRateLimiter.EXPECT().IsAllowed(gomock.Any()).Return(true)
	mockMetrics.EXPECT().IncTotalConnections()
	mockMetrics.EXPECT().IncActiveConnections()
	mockMetrics.EXPECT().DecActiveConnections()
	mockPow.EXPECT().GenerateChallenge(gomock.Any()).Return("challenge", nil)
	mockPow.EXPECT().VerifySolution(gomock.Any(), gomock.Any()).Return(true)
	mockQuotes.EXPECT().GetRandomQuote().Return("test quote")
	mockMetrics.EXPECT().IncSuccessChallenges()
	mockMetrics.EXPECT().IncTotalQuotesSent()

	// Запускаем обработку соединения в отдельной горутине
	go srv.handleConnection(context.Background(), server)

	// Имитируем клиентское взаимодействие
	fmt.Fprintf(client, "HELLO\n")
	
	// Читаем challenge
	buffer := make([]byte, 1024)
	n, err := client.Read(buffer)
	require.NoError(t, err)
	assert.Contains(t, string(buffer[:n]), "CHALLENGE")

	// Отправляем решение
	fmt.Fprintf(client, "SOLUTION test_solution\n")

	// Читаем цитату
	n, err = client.Read(buffer)
	require.NoError(t, err)
	assert.Contains(t, string(buffer[:n]), "QUOTE")
}

func TestServer_HandleConnection_RateLimit(t *testing.T) {
	srv, _, _, mockMetrics, mockRateLimiter, mockLogger := setupTest(t)

	// Создаем тестовое соединение
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()
	server = fakeTCPConn{Conn: server}

	// Настраиваем мок для имитации превышения лимита
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockRateLimiter.EXPECT().IsAllowed(gomock.Any()).Return(false)
	mockMetrics.EXPECT().IncTotalConnections()
	mockMetrics.EXPECT().IncActiveConnections()
	mockMetrics.EXPECT().DecActiveConnections()

	// Запускаем обработку соединения
	done := make(chan struct{})
	go func() {
		srv.handleConnection(context.Background(), server)
		close(done)
	}()

	// Проверяем, что соединение было закрыто
	select {
	case <-done:
		// Ожидаемое поведение
	case <-time.After(time.Second):
		t.Fatal("Connection wasn't closed after rate limit check")
	}
}

func TestServer_HandleConnection_Timeout(t *testing.T) {
	srv, _, _, mockMetrics, mockRateLimiter, mockLogger := setupTest(t)

	// Создаем пару тестовых соединений
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()
	server = fakeTCPConn{Conn: server}

	// Настройка моков
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	mockRateLimiter.EXPECT().IsAllowed(gomock.Any()).Return(true)
	mockMetrics.EXPECT().IncTotalConnections()
	mockMetrics.EXPECT().IncActiveConnections()
	mockMetrics.EXPECT().DecActiveConnections()

	// Запускаем обработку соединения
	done := make(chan struct{})
	go func() {
		srv.handleConnection(context.Background(), server)
		close(done)
	}()

	// Имитируем таймаут
	time.Sleep(2 * time.Second) // Ждем больше, чем таймаут в конфиге

	// Проверяем, что соединение было закрыто
	select {
	case <-done:
		// Ожидаемое поведение - соединение закрыто по таймауту
	case <-time.After(3 * time.Second):
		t.Fatal("Connection wasn't closed after timeout")
	}
}

func TestServer_Stop(t *testing.T) {
	srv, _, _, mockMetrics, _, mockLogger := setupTest(t)

	// Настраиваем ожидаемое поведение моков
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	// Ожидаем два вызова Wait() - один из Start, другой из Stop
	mockMetrics.EXPECT().Wait().AnyTimes()

	// Запускаем сервер
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start(ctx)
	}()

	// Даем серверу время на запуск
	time.Sleep(100 * time.Millisecond)

	// Останавливаем сервер
	cancel()
	err := srv.Stop()
	assert.NoError(t, err)

	// Ждем завершения Start
	err = <-errCh
	assert.NoError(t, err)

	// Проверяем, что listener закрыт
	_, err = net.Dial("tcp", srv.listener.Addr().String())
	assert.Error(t, err)
} 