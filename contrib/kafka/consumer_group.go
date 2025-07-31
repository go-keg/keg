package kafka

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/go-keg/keg/contrib/config"
	"github.com/go-kratos/kratos/v2/log"
)

type ConsumerGroupOption func(cg *ConsumerGroup)

type ConsumerGroupHandlerOption func(cg *ConsumerGroupHandler)

type MessageHandler func(message *sarama.ConsumerMessage) error

func SetLogger(logger log.Logger) ConsumerGroupOption {
	return func(cg *ConsumerGroup) {
		cg.logger = logger
	}
}

func SetConfigOptions(opts ...ConfigOption) ConsumerGroupOption {
	return func(cg *ConsumerGroup) {
		for _, opt := range opts {
			opt(cg.config)
		}
	}
}

func WithHandlerSetup(f func(s sarama.ConsumerGroupSession) error) ConsumerGroupHandlerOption {
	return func(h *ConsumerGroupHandler) {
		h.setup = f
	}
}

func WithHandlerCleanup(f func(s sarama.ConsumerGroupSession) error) ConsumerGroupHandlerOption {
	return func(h *ConsumerGroupHandler) {
		h.cleanup = f
	}
}

func WithLogger(logger log.Logger) ConsumerGroupHandlerOption {
	return func(h *ConsumerGroupHandler) {
		h.logger = logger
	}
}

type ConsumerGroup struct {
	topics  []string
	handler sarama.ConsumerGroupHandler
	config  *sarama.Config
	client  sarama.ConsumerGroup
	sigterm chan os.Signal
	logger  log.Logger
	log     *log.Helper
}

func NewConsumerGroup(brokers []string, group string, topics []string, handler sarama.ConsumerGroupHandler, opts ...ConsumerGroupOption) (*ConsumerGroup, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_1_0_0
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cg := &ConsumerGroup{
		config:  cfg,
		handler: handler,
		topics:  topics,
		sigterm: make(chan os.Signal, 1),
		logger:  log.DefaultLogger,
	}
	for _, opt := range opts {
		opt(cg)
	}
	client, err := sarama.NewConsumerGroup(brokers, group, cg.config)
	if err != nil {
		return nil, fmt.Errorf("error creating consumer group client: %v", err)
	}
	cg.log = log.NewHelper(log.With(cg.logger,
		"module", "kafka_consumer",
		"topics", strings.Join(topics, ","),
		"group", group,
	))
	cg.client = client
	return cg, nil
}

func NewConsumerGroupFromConfig(config config.Kafka, consumerGroup config.KafkaConsumerGroup, handler sarama.ConsumerGroupHandler, opts ...ConfigOption) (*ConsumerGroup, error) {
	opts = append(opts, SetNetSASL(config.User, config.Password))
	return NewConsumerGroup(
		config.GetAddr(),
		consumerGroup.GroupID,
		consumerGroup.Topics,
		handler,
		SetConfigOptions(opts...),
	)
}

func (r *ConsumerGroup) Run(ctx context.Context) error {
	go func() {
		for {
			err := r.client.Consume(ctx, r.topics, r.handler)
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return
			} else if err != nil {
				r.log.Errorf("error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	signal.Notify(r.sigterm, os.Interrupt, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		r.log.Debug("terminating: context cancelled")
	case <-r.sigterm:
		r.log.Debug("terminating: via signal")
	}
	return r.client.Close()
}

type ConsumerGroupHandler struct {
	setup   func(s sarama.ConsumerGroupSession) error
	cleanup func(s sarama.ConsumerGroupSession) error
	handler MessageHandler // topic -> handler
	logger  log.Logger
	log     *log.Helper
}

func NewConsumerGroupHandler(handler MessageHandler, opts ...ConsumerGroupHandlerOption) *ConsumerGroupHandler {
	h := &ConsumerGroupHandler{
		setup: func(_ sarama.ConsumerGroupSession) error {
			return nil
		},
		cleanup: func(_ sarama.ConsumerGroupSession) error {
			return nil
		},
		handler: handler,
		logger:  log.DefaultLogger,
	}
	for _, opt := range opts {
		opt(h)
	}
	h.log = log.NewHelper(log.With(h.logger, "module", "ConsumerGroupHandler"))
	return h
}

func (r *ConsumerGroupHandler) Setup(sess sarama.ConsumerGroupSession) error {
	return r.setup(sess)
}

func (r *ConsumerGroupHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	return r.cleanup(sess)
}

func (r *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				r.log.Debug("message channel was closed")
				return nil
			}
			if err := r.handler(message); err != nil {
				r.log.Errorf("message: topic = %s partition = %d offset = %d err = %v", message.Topic, message.Partition, message.Offset, err)
				claim.HighWaterMarkOffset()
				return err
			}
			r.log.Debugf("message claimed: topic = %s partition = %d offset = %d", message.Topic, message.Partition, message.Offset)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
