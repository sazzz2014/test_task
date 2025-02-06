package metrics

import "sync/atomic"

type Metrics struct {
    activeConnections   atomic.Int32
    totalConnections    atomic.Int64
    failedChallenges   atomic.Int64
    successChallenges  atomic.Int64
    totalQuotesSent    atomic.Int64
}

func NewMetrics() *Metrics {
    return &Metrics{}
}

func (m *Metrics) IncActiveConnections()    { m.activeConnections.Add(1) }
func (m *Metrics) DecActiveConnections()    { m.activeConnections.Add(-1) }
func (m *Metrics) IncTotalConnections()     { m.totalConnections.Add(1) }
func (m *Metrics) IncFailedChallenges()     { m.failedChallenges.Add(1) }
func (m *Metrics) IncSuccessChallenges()    { m.successChallenges.Add(1) }
func (m *Metrics) IncTotalQuotesSent()      { m.totalQuotesSent.Add(1) }

func (m *Metrics) GetStats() map[string]int64 {
    return map[string]int64{
        "active_connections":  int64(m.activeConnections.Load()),
        "total_connections":   m.totalConnections.Load(),
        "failed_challenges":   m.failedChallenges.Load(),
        "success_challenges": m.successChallenges.Load(),
        "total_quotes_sent":   m.totalQuotesSent.Load(),
    }
} 