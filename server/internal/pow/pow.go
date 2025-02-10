package pow

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "sync"
    "time"
)

type PoW struct {
    difficulty int
    usedSolutions sync.Map // для хранения использованных решений
    cleanupInterval time.Duration
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
    // Проверяем, не было ли уже использовано это решение
    key := challenge + solution
    if _, exists := p.usedSolutions.Load(key); exists {
        return false
    }
    
    // Проверяем решение
    combined := challenge + solution
    hash := sha256.Sum256([]byte(combined))
    
    valid := true
    for i := 0; i < p.difficulty; i++ {
        if hash[i/8]>>(7-(i%8))&1 != 0 {
            valid = false
            break
        }
    }
    
    if valid {
        // Сохраняем использованное решение
        p.usedSolutions.Store(key, time.Now())
    }
    
    return valid
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