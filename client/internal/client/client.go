package client

import (
    "bufio"
    "fmt"
    "net"
    "strings"
    "time"

    "client/internal/pow"
)

type Client struct {
    address  string
    solver   *pow.Solver
    timeout  time.Duration
}

func NewClient(address string, difficulty int, timeout time.Duration) *Client {
    return &Client{
        address:  address,
        solver:   pow.NewSolver(difficulty),
        timeout:  timeout,
    }
}

func (c *Client) GetQuote() (string, error) {
    conn, err := net.DialTimeout("tcp", c.address, c.timeout)
    if err != nil {
        return "", err
    }
    defer conn.Close()

    conn.SetDeadline(time.Now().Add(c.timeout))

    // Отправляем HELLO
    fmt.Fprintf(conn, "HELLO\n")

    reader := bufio.NewReader(conn)
    
    // Получаем вызов
    challenge, err := reader.ReadString('\n')
    if err != nil {
        return "", err
    }

    parts := strings.Fields(challenge)
    if len(parts) != 2 || parts[0] != "CHALLENGE" {
        return "", fmt.Errorf("invalid challenge format")
    }

    // Решаем PoW
    solution := c.solver.Solve(parts[1])
    
    // Отправляем решение
    fmt.Fprintf(conn, "SOLUTION %s\n", solution)

    // Получаем ответ
    response, err := reader.ReadString('\n')
    if err != nil {
        return "", err
    }

    parts = strings.Fields(response)
    if len(parts) < 1 {
        return "", fmt.Errorf("invalid response")
    }

    if parts[0] == "ERROR" {
        return "", fmt.Errorf("solution rejected")
    }

    if parts[0] != "QUOTE" {
        return "", fmt.Errorf("unexpected response")
    }

    return strings.Join(parts[1:], " "), nil
} 