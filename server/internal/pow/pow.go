package pow

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

type PoW struct {
    difficulty int
}

func NewPoW(difficulty int) *PoW {
    return &PoW{difficulty: difficulty}
}

func (p *PoW) GenerateChallenge(length int) (string, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return "", fmt.Errorf("failed to generate challenge: %w", err)
    }
    return hex.EncodeToString(bytes), nil
}

func (p *PoW) VerifySolution(challenge, solution string) bool {
    combined := challenge + solution
    hash := sha256.Sum256([]byte(combined))
    
    // Проверяем, начинается ли хеш с нужного количества нулевых битов
    for i := 0; i < p.difficulty; i++ {
        if hash[i/8]>>(7-(i%8))&1 != 0 {
            return false
        }
    }
    return true
} 