package alert

import (
	"context"
	"log"
)

type Alerter interface {
	Alert(ctx context.Context, content string) error
}

type LogAlert struct{}

func (r LogAlert) Alert(ctx context.Context, content string) error {
	log.Println(content)
	return nil
}
