package kafka

import (
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-kratos/kratos/v2/log"
)

// DLQHandler 是一个装饰器，为任何 ConsumerGroupHandler 添加重试和 DLQ 功能
type DLQHandler struct {
	wrappedHandler *ConsumerGroupHandler
	producer       sarama.SyncProducer
	dlqTopic       string
	maxRetries     int
	log            *log.Helper
}

// WrapDLQHandler 创建一个新的 DLQ 装饰器
func WrapDLQHandler(
	wrapped *ConsumerGroupHandler,
	producer sarama.SyncProducer,
	dlqTopic string,
	maxRetries int,
	logger log.Logger,
) sarama.ConsumerGroupHandler {
	return &DLQHandler{
		wrappedHandler: wrapped,
		producer:       producer,
		dlqTopic:       dlqTopic,
		maxRetries:     maxRetries,
		log:            log.NewHelper(log.With(logger, "module", "dlq_handler")),
	}
}

// Setup 和 Cleanup 直接调用被包裹的 handler 的方法
func (h *DLQHandler) Setup(session sarama.ConsumerGroupSession) error {
	return h.wrappedHandler.Setup(session)
}

func (h *DLQHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return h.wrappedHandler.Cleanup(session)
}

// ConsumeClaim 实现 DLQ 逻辑
func (h *DLQHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		topic := strings.TrimSuffix(message.Topic, h.dlqTopic)
		var err error
		for i := 0; i < h.maxRetries+1; i++ {
			err = h.wrappedHandler.handler(message)
			if err == nil {
				break
			}
			h.log.Warnf("message processing failed, preparing to retry: topic=%s, offset=%d, retry=%d, err=%v", message.Topic, message.Offset, i, err)
		}

		if err != nil {
			// 所有重试后仍然失败，发送到 DLQ
			h.log.Errorf("all retries failed, sent to the dead letter queue: topic=%s, offset=%d, err=%v", message.Topic, message.Offset, err)

			dlqMessage := &sarama.ProducerMessage{
				Topic: topic + h.dlqTopic,
				Key:   sarama.ByteEncoder(message.Key),
				Value: sarama.ByteEncoder(message.Value),
				Headers: []sarama.RecordHeader{
					{Key: []byte("original_topic"), Value: []byte(topic)},
					{Key: []byte("dlq_reason"), Value: []byte(err.Error())},
					{Key: []byte("failed_at"), Value: []byte(time.Now().Format(time.RFC3339))},
				},
			}

			if _, _, err = h.producer.SendMessage(dlqMessage); err != nil {
				h.log.Errorf("failed to send to the dead letter queue: %v", err)
			}
			time.Sleep(time.Second)
		}

		session.MarkMessage(message, "")
	}
	return nil
}
