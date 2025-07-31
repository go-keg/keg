package kafka

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
	klog "github.com/go-kratos/kratos/v2/log"
)

var logger = klog.DefaultLogger

var handler = NewConsumerGroupHandler(func(message *sarama.ConsumerMessage) error {
	switch message.Topic {
	case "test":
		if message.Offset%5 == 4 {
			return fmt.Errorf("offset is %d", message.Offset)
		}
		fmt.Println("consume", message.Topic, message.Partition, message.Offset)
	}
	return nil
}, WithLogger(logger))

func TestNewConsumerGroup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	wg := sync.WaitGroup{}

	cg, err := NewConsumerGroup(
		[]string{"localhost:9092"},
		fmt.Sprintf("test-%s", time.Now().Format("15-04-05")),
		[]string{"test"},
		handler,
		SetLogger(logger),
	)
	if err != nil {
		t.Fatal(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := cg.Run(ctx)
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

func TestNewConsumerGroupDLQ(t *testing.T) {
	producer, err := NewSyncProducer([]string{"localhost:9092"})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.TODO())
	wg := sync.WaitGroup{}

	cg, err := NewConsumerGroup(
		[]string{"localhost:9092"},
		fmt.Sprintf("test-%s", time.Now().Format("2006-01-02_15-04-05")),
		[]string{"test", "test_dlq"},
		WrapDLQHandler(handler, producer, "_dlq", 3, logger),
		SetLogger(logger),
	)
	if err != nil {
		t.Fatal(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = cg.Run(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// 启动10s后主动取消，终止运行
	go func() {
		<-time.Tick(60 * time.Second)
		cancel()
	}()
	wg.Wait()
}
