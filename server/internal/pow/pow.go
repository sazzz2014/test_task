package pow

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "sync"
    "time"
    "sync/atomic"
)

type PoW struct {
    difficulty int
    usedSolutions sync.Map // для хранения использованных решений
    cleanupInterval time.Duration
    totalAttempts   atomic.Int64
    validSolutions  atomic.Int64
    replayAttempts  atomic.Int64
}

func NewPoW(difficulty int) *PoW {
    pow := &PoW{
        difficulty: difficulty,
        cleanupInterval: 5 * time.Minute,
    }
    go pow.cleanupRoutine()
    return pow
}

func (p *PoW) GenerateChallenge(length int) (string, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return "", fmt.Errorf("failed to generate challenge: %w", err)
    }
    return hex.EncodeToString(bytes), nil
}

func (p *PoW) VerifySolution(challenge, solution string) bool {
    p.totalAttempts.Add(1)
    
    // Проверка на replay
    key := challenge + solution
    if _, exists := p.usedSolutions.Load(key); exists {
        p.replayAttempts.Add(1)
        return false
    }
    
    // Проверка формата solution
    if _, err := hex.DecodeString(solution); err != nil {
        return false
    }
    
    // Проверка длины solution
    if len(solution) > 64 { 
        return false
    }
    
    // Проверяем решение
    combined := challenge + solution
    hash := sha256.Sum256([]byte(combined))
    
    // Оптимизированная битовая проверка
    for i := 0; i < p.difficulty; i++ {
        if (hash[i/8] >> (7 - (i % 8))) & 1 != 0 {
            return false
        }
    }
    
    // Сохраняем решение с временем жизни
    p.validSolutions.Add(1)
    p.usedSolutions.Store(key, time.Now())
    
    return true
}

func (p *PoW) cleanupRoutine() {
    ticker := time.NewTicker(p.cleanupInterval)
    for range ticker.C {
        now := time.Now()
        p.usedSolutions.Range(func(key, value interface{}) bool {
            if timestamp, ok := value.(time.Time); ok {
                if now.Sub(timestamp) > p.cleanupInterval {
                    p.usedSolutions.Delete(key)
                }
            }
            return true
        })
    }
}

func (p *PoW) GenerateAndVerify(challenge, solution string) (bool, error) {
    // Проверка решения
    return p.VerifySolution(challenge, solution), nil
}

func (p *PoW) GetStats() map[string]int64 {
    return map[string]int64{
        "total_attempts":   p.totalAttempts.Load(),
        "valid_solutions": p.validSolutions.Load(),
        "replay_attempts": p.replayAttempts.Load(),
    }
} 