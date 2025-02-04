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
    "server/internal/pow"
    "server/internal/protocol"
    "server/internal/quotes"
)

type Server struct {
    config       *config.Config
    pow          *pow.PoW
    quoteService *quotes.QuoteService
    listener     net.Listener
    connections  sync.Map
    shutdown     chan struct{}
}

func NewServer(cfg *config.Config, pow *pow.PoW, qs *quotes.QuoteService) *Server {
    return &Server{
        config:       cfg,
        pow:          pow,
        quoteService: qs,
        shutdown:     make(chan struct{}),
    }
}

func (s *Server) Start(ctx context.Context) error {
    var err error
    s.listener, err = net.Listen("tcp", s.config.Port)
    if err != nil {
        return err
    }

    go s.acceptConnections(ctx)

    <-ctx.Done()
    return s.Stop()
}

func (s *Server) Stop() error {
    close(s.shutdown)
    if s.listener != nil {
        return s.listener.Close()
    }
    return nil
}

func (s *Server) acceptConnections(ctx context.Context) {
    var connectionCount int32
    for {
        conn, err := s.listener.Accept()
        if err != nil {
            select {
            case <-s.shutdown:
                return
            default:
                continue
            }
        }

        if connectionCount >= int32(s.config.MaxConnections) {
            conn.Close()
            continue
        }

        connectionCount++
        go s.handleConnection(ctx, conn, &connectionCount)
    }
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn, count *int32) {
    defer func() {
        conn.Close()
        *count--
    }()

    s.connections.Store(conn.RemoteAddr(), conn)
    defer s.connections.Delete(conn.RemoteAddr())

    reader := bufio.NewReader(conn)
    
    // Установка таймаутов
    conn.SetReadDeadline(time.Now().Add(time.Duration(s.config.ReadTimeout) * time.Second))
    conn.SetWriteDeadline(time.Now().Add(time.Duration(s.config.WriteTimeout) * time.Second))

    // Ожидаем HELLO
    message, err := reader.ReadString('\n')
    if err != nil || !strings.HasPrefix(message, protocol.CmdHello) {
        return
    }

    // Генерируем и отправляем вызов
    challenge := s.pow.GenerateChallenge(s.config.ChallengeLength)
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
        quote := s.quoteService.GetRandomQuote()
        fmt.Fprintf(conn, "%s %s\n", protocol.CmdQuote, quote)
    } else {
        fmt.Fprintf(conn, "%s\n", protocol.CmdError)
    }
} 