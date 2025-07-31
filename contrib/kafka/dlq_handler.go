package kafka

import (
	"fmt"
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

// ConsumeClaim 是核心，在这里实现 DLQ 逻辑
func (h *DLQHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		topic := strings.TrimSuffix(message.Topic, h.dlqTopic)
		handler, handlerExists := h.wrappedHandler.handlers[topic]
		if !handlerExists {
			return fmt.Errorf("topic: %s, not found hand", topic)
		}

		var err error
		for i := 0; i < h.maxRetries+1; i++ {
			err = handler(message)
			if err == nil {
				break // 处理成功，跳出重试循环
			}
			h.log.Warnf("消息处理失败，准备重试: topic=%s, offset=%d, retry=%d, err=%v", message.Topic, message.Offset, i, err)
		}

		if err != nil {
			// 所有重试后仍然失败，发送到 DLQ
			h.log.Errorf("所有重试失败，发送到死信队列: topic=%s, offset=%d, err=%v", message.Topic, message.Offset, err)

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
				h.log.Errorf("发送到死信队列失败: %v", err)
			}
			time.Sleep(time.Second)
		}

		// 关键：无论成功、失败进入DLQ，都必须标记原始消息为已消费
		session.MarkMessage(message, "")
	}
	return nil
}
