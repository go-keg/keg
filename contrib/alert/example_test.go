package alert

import (
	"fmt"
	"testing"
	"time"
)

func TestDeduper(t *testing.T) {
	d := NewDeduper(SetTTL(2*time.Minute), SetThreshold(0.90), SetCleanupInterval(2*time.Second), SetReport(func(record *ErrorRecord) {
		if record.Count == 1 {
			fmt.Println(record.RawMsg)
		}
	}))
	msgs := []string{
		"Error: failed to fetch user id=12345",
		"Error: failed to fetch user id=67890",
		"Error: failed to fetch user id=14725",
		"Error: failed to connect to db at 10.0.0.1",
		"Error: failed to connect to db at 10.0.0.2",
		"Error: failed to connect to db at 10.0.0.3",
		`{"level":"error","service.id":"golang-shop-service","service.name":"shop-service","service.version":"v1.6.0","ts":"2025-11-13T09:23:43Z","caller":"job/ga_custom_channel_group_metrics.go:33","module":"job","method":"createChannelGroupMetrics","shopDomain":"dji.myshopify.com","err":"googleapi: Error 403: The caller does not have permission, forbidden"}`,
		`{"level":"error","service.id":"golang-shop-service","service.name":"shop-service","service.version":"v1.6.0","ts":"2025-11-13T09:23:43Z","caller":"job/ga_custom_channel_group_metrics.go:33","module":"job","method":"createChannelGroupMetrics","shopDomain":"xiaomi.myshopify.com","err":"googleapi: Error 403: The caller does not have permission, forbidden"}`,
		`{"level":"error","service.id":"golang-shop-service","service.name":"shop-service","service.version":"v1.6.0","ts":"2025-11-13T09:23:43Z","caller":"job/ga_custom_channel_group_metrics.go:33","module":"job","method":"createChannelGroupMetrics","shopDomain":"huawei.myshopify.com","err":"googleapi: Error 403: The caller does not have permission, forbidden"}`,
	}
	for i := 0; i < 100; i++ {
		for _, m := range msgs {
			d.Record(m)
			time.Sleep(200 * time.Millisecond)
		}
	}
}
