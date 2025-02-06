package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"server/internal/logger"
	"server/internal/protocol"
)

type protocolHandler struct {
	conn   net.Conn
	reader *bufio.Reader
	logger *logger.Logger
}

func (h *protocolHandler) handleHello() error {
	message, err := h.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read HELLO: %w", err)
	}
	if !strings.HasPrefix(message, protocol.CmdHello) {
		return fmt.Errorf("invalid hello message")
	}
	return nil
}

func (h *protocolHandler) handleChallenge(challenge string) error {
	_, err := fmt.Fprintf(h.conn, "%s %s\n", protocol.CmdChallenge, challenge)
	return err
}

// ... другие методы обработки протокола 