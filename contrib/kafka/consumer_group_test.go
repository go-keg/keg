package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	klog "github.com/go-kratos/kratos/v2/log"
	"log"
	"sync"
	"testing"
	"time"
)

func TestNewConsumerGroup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	wg := sync.WaitGroup{}

	cg := NewConsumerGroup(
		[]string{"127.0.0.1:9093"},
		fmt.Sprintf("test-%s", time.Now().Format("15-04-05")),
		[]string{"test"},
		SetLogger(klog.DefaultLogger),
		SetConfigOptions(
			SetNetSASL("", ""),
			SetConsumerGroupBalanceStrategy(sarama.NewBalanceStrategyRange()),
			// ...
		),
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := cg.Run(ctx, func(message *sarama.ConsumerMessage) error {
			log.Printf("todo: topic=%s partition=%d offset=%d", message.Topic, message.Partition, message.Offset)
			time.Sleep(time.Millisecond * 500)
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}()

	// 启动10s后主动取消，终止运行
	go func() {
		<-time.Tick(10 * time.Second)
		cancel()
	}()
	wg.Wait()
}
