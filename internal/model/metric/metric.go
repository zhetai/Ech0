package model

import (
	"time"
)

// CpuMetric cpu监控指标
type CpuMetric struct {
	UsagePercent float64 // CPU 使用率百分比
	Cores        int     // CPU 核心数
	FrequencyMHz uint64  // CPU 主频，单位 MHz
}

// MemoryMetric 内存监控指标
type MemoryMetric struct {
	Total      uint64  // 总内存大小
	Used       uint64  // 已使用内存大小
	Available  uint64  // 可用内存大小
	Percentage float64 // 内存使用率百分比
}

// DiskMetric 磁盘监控指标
type DiskMetric struct {
	Total      uint64  // 总磁盘大小
	Used       uint64  // 已使用磁盘大小
	Available  uint64  // 可用磁盘大小
	Percentage float64 // 磁盘使用率百分比
}

// NetworkMetric 网络监控指标
type NetworkMetric struct {
	TotalBytesSent         uint64  // 总发送字节数
	TotalBytesReceived     uint64  // 总接收字节数
	BytesSentPerSecond     float64 // 每秒发送字节数 (B/s)
	BytesReceivedPerSecond float64 // 每秒接收字节数 (B/s)
}

// SystemMetric 系统监控指标
type SystemMetric struct {
	Hostname      string        // 主机名
	OsName        string        // 操作系统名称
	Uptime        time.Duration // 系统运行时长
	KernelVersion string        // 内核版本
	KernelArch    string        // 内核架构
	Time          time.Time     // 采样时间
	TimeZone      string        // 采样时区
	ProcessCount int           // 当前进程数
	ThreadCount  int           // 当前线程数
	GolangVersion  string        // Golang 版本
	GoRoutineCount int           // 当前 Goroutine 数量
}

// Metrics 综合监控指标
type Metrics struct {
	CPU     CpuMetric     // CPU 监控指标
	Memory  MemoryMetric  // 内存监控指标
	Disk    DiskMetric    // 磁盘监控指标
	Network NetworkMetric // 网络监控指标
	System  SystemMetric  // 系统监控指标
}
