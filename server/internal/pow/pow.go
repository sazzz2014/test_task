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
    leadingZeros := 0
    for i := 0; i < len(hash); i++ {
        if hash[i] == 0 {
            leadingZeros += 8
        } else {
            leadingZeros += countLeadingZeros(hash[i])
            break
        }
    }
    return leadingZeros >= p.difficulty
}

func countLeadingZeros(b byte) int {
    count := 0
    for b&0x80 == 0 {
        count++
        b <<= 1
    }
    return count
}

func (p *PoW) GenerateAndVerify(challenge, solution string) (bool, error) {
    // Проверка решения
    return p.VerifySolution(challenge, solution), nil
} 