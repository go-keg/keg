package alert

import (
	"testing"
	"time"
)

func TestDeduper(t *testing.T) {
	d := NewDeduper(
		// SetAlert(workwechat.NewWebhook(os.Getenv("WORKWECHAT_WEBHOOK_TOKEN"))),
		SetReport(FibReport),
		SetTTL(5*time.Minute),
	)
	msgs := []string{
		"Error: failed to fetch user id=12345",
		"Error: failed to fetch user id=67890",
		"Error: failed to fetch user id=14725",
		"Error: failed to connect to db at 10.0.0.1",
		"Error: failed to connect to db at 10.0.0.2",
		"Error: failed to connect to db at 10.0.0.3",
		`{"level":"error","service.id":"golang-shop-service","service.name":"shop-service","service.version":"v1.6.0","ts":"2025-11-13T09:23:43Z","caller":"job/ga_metrics.go:33","shopDomain":"dji.myshopify.com","err":"googleapi: Error 403: The caller does not have permission, forbidden"}`,
		`{"level":"error","service.id":"golang-shop-service","service.name":"shop-service","service.version":"v1.6.0","ts":"2025-11-13T09:23:43Z","caller":"job/ga_metrics.go:33","shopDomain":"xiaomi.myshopify.com","err":"googleapi: Error 403: The caller does not have permission, forbidden"}`,
		`{"level":"error","service.id":"golang-shop-service","service.name":"shop-service","service.version":"v1.6.0","ts":"2025-11-13T09:23:43Z","caller":"job/ga_metrics.go:33","shopDomain":"huawei.myshopify.com","err":"googleapi: Error 403: The caller does not have permission, forbidden"}`,
	}
	for i := 0; i < 100; i++ {
		for _, m := range msgs {
			_ = d.Alert(t.Context(), m)
			time.Sleep(200 * time.Millisecond)
		}
	}
}
