package pow

import (
    "crypto/sha256"
    "encoding/hex"
    "math/rand"
    "time"
)

type PoW struct {
    difficulty int
}

func NewPoW(difficulty int) *PoW {
    return &PoW{difficulty: difficulty}
}

func (p *PoW) GenerateChallenge(length int) string {
    rand.Seed(time.Now().UnixNano())
    bytes := make([]byte, length)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
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