package alert

import (
	"context"
	"sync"
	"time"
)

type ErrorRecord struct {
	RawMsg     string    // 原始错误信息
	Normalized string    // 归一化后的内容
	Count      int       // 计数
	AlertCount int       // 告警计数
	FirstSeen  time.Time // 首次出现时间
	LastSeen   time.Time // 最后出现时间
}

type Deduper struct {
	mu              sync.Mutex
	cache           []*ErrorRecord
	threshold       float64       // 相似度阈值（0~1）
	ttl             time.Duration // 记录保留时间
	cleanupInterval time.Duration
	report          Report
	alert           Alerter
}

type DeduperOption func(deduper *Deduper)

func SetThreshold(val float64) DeduperOption {
	return func(deduper *Deduper) {
		deduper.threshold = val
	}
}

func SetTTL(val time.Duration) DeduperOption {
	return func(deduper *Deduper) {
		deduper.ttl = val
	}
}

func SetCleanupInterval(val time.Duration) DeduperOption {
	return func(deduper *Deduper) {
		deduper.cleanupInterval = val
	}
}

func SetReport(val Report) DeduperOption {
	return func(deduper *Deduper) {
		deduper.report = val
	}
}

func SetAlert(val Alerter) DeduperOption {
	return func(deduper *Deduper) {
		deduper.alert = val
	}
}

func NewDeduper(opts ...DeduperOption) *Deduper {
	d := &Deduper{
		cache:     []*ErrorRecord{},
		threshold: 0.9,
		ttl:       time.Hour * 24,
		report:    DefaultReport,
		alert:     LogAlert{},
	}
	for _, opt := range opts {
		opt(d)
	}
	if d.cleanupInterval == 0 {
		d.cleanupInterval = d.ttl / 2
	}
	go d.cleanupLoop()
	return d
}

func (d *Deduper) Alert(ctx context.Context, msg string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	norm := NormalizeError(msg)
	now := time.Now()

	for _, rec := range d.cache {
		if Similar(norm, rec.Normalized) >= d.threshold {
			rec.Count++
			rec.LastSeen = now
			if content, ok := d.report(ctx, rec); ok {
				return d.alert.Alert(ctx, content)
			}
			return nil
		}
	}

	rec := &ErrorRecord{
		RawMsg:     msg,
		Normalized: norm,
		Count:      1,
		FirstSeen:  now,
		LastSeen:   now,
	}
	d.cache = append(d.cache, rec)
	if content, ok := d.report(ctx, rec); ok {
		return d.alert.Alert(ctx, content)
	}
	return nil
}

// 定期清理过期记录
func (d *Deduper) cleanupLoop() {
	ticker := time.NewTicker(d.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		d.mu.Lock()
		now := time.Now()
		newCache := make([]*ErrorRecord, 0, len(d.cache))
		for _, rec := range d.cache {
			if now.Sub(rec.FirstSeen) < d.ttl {
				newCache = append(newCache, rec)
			}
		}
		d.cache = newCache
		d.mu.Unlock()
	}
}
