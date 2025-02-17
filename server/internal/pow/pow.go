package pow

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type PoW struct {
	difficulty      int           // Сложность - количество начальных нулевых битов
	usedSolutions   sync.Map      // Thread-safe карта для хранения использованных решений
	cleanupInterval time.Duration // Интервал очистки старых решений
	solutionTTL     time.Duration // Время жизни решения
	totalAttempts   atomic.Int64  // Счетчик всех попыток
	validSolutions  atomic.Int64  // Счетчик успешных решений
	replayAttempts  atomic.Int64  // Счетчик попыток повторного использования
}

func NewPoW(difficulty int) *PoW {
	pow := &PoW{
		difficulty:      difficulty,
		cleanupInterval: 5 * time.Minute,
		solutionTTL:     10 * time.Minute, // Устанавливаем TTL
	}
	go pow.cleanupRoutine()
	return pow
}

// Генерирует случайный challenge заданной длины
func (p *PoW) GenerateChallenge(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate challenge: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// Проверяет решение от клиента
func (p *PoW) VerifySolution(challenge, solution string) bool {
	p.totalAttempts.Add(1)

	// 1. Базовые проверки
	if len(challenge) == 0 || len(solution) == 0 {
		return false
	}

	// 2. Защита от replay-атак
	key := challenge + solution
	if _, exists := p.usedSolutions.Load(key); exists {
		p.replayAttempts.Add(1)
		return false
	}

	// 3. Проверка формата solution
	if _, err := hex.DecodeString(solution); err != nil {
		return false
	}

	// 4. Проверка длины solution
	if len(solution) > 64 {
		return false
	}

	// 5. Проверка решения (битовая проверка)
	combined := challenge + solution
	hash := sha256.Sum256([]byte(combined))

	// Оптимизированная битовая проверка
	mask := byte(0xFF) << (8 - (p.difficulty % 8))
	for i := 0; i < p.difficulty/8; i++ {
		if hash[i] != 0 {
			return false
		}
	}
	if p.difficulty%8 > 0 {
		if (hash[p.difficulty/8] & mask) != 0 {
			return false
		}
	}

	// 6. Сохранение успешного решения
	p.validSolutions.Add(1)
	p.usedSolutions.Store(key, time.Now())

	return true
}

// Периодически очищает старые решения
func (p *PoW) cleanupRoutine() {
	ticker := time.NewTicker(p.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		p.usedSolutions.Range(func(key, value interface{}) bool {
			if timestamp, ok := value.(time.Time); ok {
				if now.Sub(timestamp) > p.solutionTTL {
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

// Возвращает статистику работы
func (p *PoW) GetStats() map[string]int64 {
	return map[string]int64{
		"total_attempts":  p.totalAttempts.Load(),
		"valid_solutions": p.validSolutions.Load(),
		"replay_attempts": p.replayAttempts.Load(),
	}
}

func (p *PoW) GetDetailedStats() map[string]interface{} {
	return map[string]interface{}{
		"total_attempts":           p.totalAttempts.Load(),
		"valid_solutions":          p.validSolutions.Load(),
		"replay_attempts":          p.replayAttempts.Load(),
		"difficulty":               p.difficulty,
		"solution_ttl_minutes":     p.solutionTTL.Minutes(),
		"cleanup_interval_minutes": p.cleanupInterval.Minutes(),
	}
}
