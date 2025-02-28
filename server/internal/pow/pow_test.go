package pow

import (
    "testing"
    "time"
    "encoding/hex"
    "strings"
    "sync"
)

func TestPoW_GenerateChallenge(t *testing.T) {
    pow := NewPoW(4)
    
    // Проверяем генерацию challenge разной длины
    lengths := []int{8, 16, 32}
    for _, length := range lengths {
        challenge, err := pow.GenerateChallenge(length)
        if err != nil {
            t.Errorf("Failed to generate challenge: %v", err)
        }
        
        // Проверяем длину
        if len(challenge) != length*2 { // *2 так как hex encoding
            t.Errorf("Expected length %d, got %d", length*2, len(challenge))
        }
        
        // Проверяем, что это валидный hex
        _, err = hex.DecodeString(challenge)
        if err != nil {
            t.Errorf("Generated challenge is not valid hex: %v", err)
        }
    }
    
    // Проверяем уникальность
    challenges := make(map[string]bool)
    for i := 0; i < 100; i++ {
        challenge, _ := pow.GenerateChallenge(16)
        if challenges[challenge] {
            t.Error("Generated duplicate challenge")
        }
        challenges[challenge] = true
    }
}

func TestPoW_VerifySolution(t *testing.T) {
    pow := NewPoW(4)
    
    tests := []struct {
        name      string
        challenge string
        solution  string
        want      bool
    }{
        {
            name:      "Empty inputs",
            challenge: "",
            solution:  "",
            want:      false,
        },
        {
            name:      "Invalid hex solution",
            challenge: "1234",
            solution:  "xyz",
            want:      false,
        },
        {
            name:      "Solution too long",
            challenge: "1234",
            solution:  strings.Repeat("0", 65),
            want:      false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := pow.VerifySolution(tt.challenge, tt.solution); got != tt.want {
                t.Errorf("VerifySolution() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestPoW_ReplayAttack(t *testing.T) {
    pow := NewPoW(4)
    challenge := "1234"
    solution := "5678"
    
    // Первая попытка должна пройти
    if !pow.VerifySolution(challenge, solution) {
        t.Error("First attempt should succeed")
    }
    
    // Повторная попытка должна быть отклонена
    if pow.VerifySolution(challenge, solution) {
        t.Error("Replay attack should be prevented")
    }
}

func TestPoW_Concurrency(t *testing.T) {
    pow := NewPoW(4)
    var wg sync.WaitGroup
    
    // Запускаем множество горутин
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Add(-1)
            challenge, _ := pow.GenerateChallenge(16)
            pow.VerifySolution(challenge, "test")
        }()
    }
    
    wg.Wait()
}

func TestPoW_CleanupRoutine(t *testing.T) {
    pow := &PoW{
        difficulty:      4,
        cleanupInterval: 100 * time.Millisecond,
        solutionTTL:    200 * time.Millisecond,
    }
    
    // Добавляем тестовое решение
    pow.usedSolutions.Store("test", time.Now())
    
    // Ждем, пока пройдет TTL
    time.Sleep(300 * time.Millisecond)
    
    // Проверяем, что решение было удалено
    _, exists := pow.usedSolutions.Load("test")
    if exists {
        t.Error("Solution should be cleaned up")
    }
}

func TestPoW_Difficulty(t *testing.T) {
    // Тестируем разные уровни сложности
    difficulties := []int{1, 4, 8}
    
    for _, diff := range difficulties {
        t.Run(fmt.Sprintf("difficulty_%d", diff), func(t *testing.T) {
            pow := NewPoW(diff)
            challenge, _ := pow.GenerateChallenge(16)
            
            // Генерируем решение (в реальном приложении это делает клиент)
            solution := findSolution(challenge, diff)
            
            if !pow.VerifySolution(challenge, solution) {
                t.Error("Valid solution was rejected")
            }
        })
    }
}

// Вспомогательная функция для поиска решения
func findSolution(challenge string, difficulty int) string {
    nonce := 0
    for {
        solution := fmt.Sprintf("%x", nonce)
        combined := challenge + solution
        hash := sha256.Sum256([]byte(combined))
        
        valid := true
        for i := 0; i < difficulty; i++ {
            if (hash[i/8] >> (7 - (i % 8))) & 1 != 0 {
                valid = false
                break
            }
        }
        
        if valid {
            return solution
        }
        nonce++
    }
} 