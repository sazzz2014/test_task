package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

// Solver реализует решатель Proof of Work
type Solver struct {
	// difficulty определяет сложность решения (количество начальных нулевых битов)
	difficulty int
}

// NewSolver создает новый экземпляр решателя с заданной сложностью
func NewSolver(difficulty int) *Solver {
	return &Solver{difficulty: difficulty}
}

// Solve находит решение для данного вызова
// challenge - строка вызова от сервера
// Возвращает найденное решение в виде hex строки
func (s *Solver) Solve(challenge string) string {
	rand.Seed(time.Now().UnixNano())
	
	for {
		nonce := make([]byte, 8)
		rand.Read(nonce)
		solution := hex.EncodeToString(nonce)
		
		combined := challenge + solution
		hash := sha256.Sum256([]byte(combined))
		
		valid := true
		for i := 0; i < s.difficulty; i++ {
			if hash[i/8]>>(7-(i%8))&1 != 0 {
				valid = false
				break
			}
		}
		
		if valid {
			return solution
		}
	}
} 