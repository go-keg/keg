package kafka

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/go-keg/keg/contrib/config"
	"github.com/go-kratos/kratos/v2/log"
)

type ConsumerGroupOption func(cg *ConsumerGroup)

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

func SetHandlerSetup(h func(s sarama.ConsumerGroupSession) error) ConsumerGroupOption {
	return func(cg *ConsumerGroup) {
		cg.setup = h
	}
}

func SetHandlerCleanup(h func(s sarama.ConsumerGroupSession) error) ConsumerGroupOption {
	return func(cg *ConsumerGroup) {
		cg.cleanup = h
	}
}

type ConsumerGroup struct {
	topics  []string
	setup   func(s sarama.ConsumerGroupSession) error
	cleanup func(s sarama.ConsumerGroupSession) error
	handler func(message *sarama.ConsumerMessage) error
	config  *sarama.Config
	client  sarama.ConsumerGroup
	sigterm chan os.Signal
	logger  log.Logger
	log     *log.Helper
	run     bool
}

func NewConsumerGroup(brokers []string, group string, topics []string, opts ...ConsumerGroupOption) *ConsumerGroup {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_1_0_0
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cg := &ConsumerGroup{
		config: cfg,
		setup: func(s sarama.ConsumerGroupSession) error {
			return nil
		},
		cleanup: func(s sarama.ConsumerGroupSession) error {
			return nil
		},
		topics:  topics,
		sigterm: make(chan os.Signal, 1),
		logger:  log.DefaultLogger,
	}
	for _, opt := range opts {
		opt(cg)
	}
	client, err := sarama.NewConsumerGroup(brokers, group, cg.config)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v", err)
	}
	cg.log = log.NewHelper(log.With(cg.logger,
		"module", "kafka_consumer",
		"topics", strings.Join(topics, ","),
		"group", group,
	))
	cg.client = client
	return cg
}

func NewConsumerGroupFromConfig(config config.Kafka, consumerGroup config.KafkaConsumerGroup, opts ...ConfigOption) *ConsumerGroup {
	opts = append(opts, SetNetSASL(config.User, config.Password))
	return NewConsumerGroup(
		config.GetAddr(),
		consumerGroup.GroupID,
		consumerGroup.Topics,
		SetConfigOptions(opts...),
	)
}

func (r *ConsumerGroup) Setup(sess sarama.ConsumerGroupSession) error {
	return r.setup(sess)
}

func (r *ConsumerGroup) Cleanup(sess sarama.ConsumerGroupSession) error {
	return r.cleanup(sess)
}

func (r *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				r.log.Debug("message channel was closed")
				return nil
			}
			if err := r.handler(message); err != nil {
				r.log.Debugf("message error: topic = %s partition = %d offset = %d err = %v", message.Topic, message.Partition, message.Offset, err)
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

func (r *ConsumerGroup) Run(ctx context.Context, handler func(message *sarama.ConsumerMessage) error) error {
	r.handler = handler
	r.run = true
	go func() {
		for {
			if r.run {
				err := r.client.Consume(ctx, r.topics, r)
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				} else if err != nil {
					r.log.Errorf("error from consumer: %v", err)
				}
				if ctx.Err() != nil {
					return
				}
			}
		}
	}()

	signal.Notify(r.sigterm, os.Interrupt, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		r.run = false
		r.log.Debug("terminating: context cancelled")
	case <-r.sigterm:
		r.run = false
		r.log.Debug("terminating: via signal")
	}
	return r.client.Close()
}
