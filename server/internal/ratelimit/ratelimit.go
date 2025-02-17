package ratelimit

import (
	"sync"
	"time"
	
	"server/internal/config"
)

// IPControl управляет ограничением запросов по IP
type IPControl struct {
	// mu мьютекс для синхронизации доступа к данным
	mu            sync.RWMutex
	// requests хранит историю запросов для каждого IP
	requests      map[string][]time.Time
	// blacklist содержит заблокированные IP и время окончания блокировки
	blacklist     map[string]time.Time
	// config содержит настройки ограничений
	config        *config.Config
}

// NewIPControl создает новый экземпляр контроля IP
func NewIPControl(cfg *config.Config) *IPControl {
	ic := &IPControl{
		requests:  make(map[string][]time.Time),
		blacklist: make(map[string]time.Time),
		config:    cfg,
	}
	
	// Периодическая очистка старых записей
	go ic.cleanup()
	
	return ic
}

// IsAllowed проверяет, разрешен ли доступ для данного IP
func (ic *IPControl) IsAllowed(ip string) bool {
	ic.mu.Lock()
	defer ic.mu.Unlock()
	
	// Проверка черного списка
	if bannedUntil, exists := ic.blacklist[ip]; exists {
		if time.Now().Before(bannedUntil) {
			return false
		}
		delete(ic.blacklist, ip)
	}
	
	// Очистка старых запросов
	now := time.Now()
	window := now.Add(-ic.config.RateLimitWindow)
	
	requests := ic.requests[ip]
	validRequests := requests[:0]
	
	for _, t := range requests {
		if t.After(window) {
			validRequests = append(validRequests, t)
		}
	}
	
	ic.requests[ip] = validRequests
	
	// Проверка лимита
	if len(validRequests) >= ic.config.MaxRequestsPerIP {
		ic.blacklist[ip] = now.Add(ic.config.BlacklistDuration)
		return false
	}
	
	// Добавление нового запроса
	ic.requests[ip] = append(ic.requests[ip], now)
	return true
}

// cleanup периодически очищает устаревшие записи
func (ic *IPControl) cleanup() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		ic.mu.Lock()
		now := time.Now()
		
		// Очистка черного списка
		for ip, bannedUntil := range ic.blacklist {
			if now.After(bannedUntil) {
				delete(ic.blacklist, ip)
			}
		}
		
		// Очистка старых запросов
		window := now.Add(-ic.config.RateLimitWindow)
		for ip, requests := range ic.requests {
			validRequests := requests[:0]
			for _, t := range requests {
				if t.After(window) {
					validRequests = append(validRequests, t)
				}
			}
			if len(validRequests) == 0 {
				delete(ic.requests, ip)
			} else {
				ic.requests[ip] = validRequests
			}
		}
		
		ic.mu.Unlock()
	}
} 