package monitor

import (
	"log"
	"sync"
	"time"

	"github.com/lin-snow/ech0/internal/metric"
	model "github.com/lin-snow/ech0/internal/model/metric"
)

type Monitor struct {
	collector metric.MetricCollector
	metrics   model.Metrics
	mu        sync.RWMutex
	interval  time.Duration
	stopChan  chan struct{}
	running   bool
}

// NewMonitor 创建一个新的监控器，指定采样器与采样周期。
func NewMonitor(collector metric.MetricCollector) *Monitor {
	return &Monitor{
		collector: collector,
		interval:  30 * time.Second,
		stopChan:  make(chan struct{}),
	}
}

// Start 开始定时采集系统指标。
func (m *Monitor) Start() {
	if m.running {
		return
	}
	m.running = true

	go func() {
		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := m.collect(); err != nil {
					log.Printf("[Monitor] collect error: %v", err)
				}
			case <-m.stopChan:
				log.Println("[Monitor] stopped")
				m.running = false
				return
			}
		}
	}()
}

// Stop 停止监控。
func (m *Monitor) Stop() {
	if !m.running {
		return
	}
	close(m.stopChan)
}

// collect 内部采样逻辑。
func (m *Monitor) collect() error {
	data, err := m.collector.Collect()
	if err != nil {
		return err
	}

	m.mu.Lock()
	m.metrics = data
	m.mu.Unlock()
	return nil
}

// GetMetrics 获取当前缓存的最新指标。
func (m *Monitor) GetMetrics() model.Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.metrics
}

// ForceCollect 立即采样一次（同步）。
func (m *Monitor) ForceCollect() (model.Metrics, error) {
	data, err := m.collector.Collect()
	if err != nil {
		return data, err
	}

	m.mu.Lock()
	m.metrics = data
	m.mu.Unlock()
	return data, nil
}
