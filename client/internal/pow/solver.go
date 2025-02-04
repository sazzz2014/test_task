package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

type Solver struct {
	difficulty int
}

func NewSolver(difficulty int) *Solver {
	return &Solver{difficulty: difficulty}
}

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