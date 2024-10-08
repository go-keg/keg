package kafka

import (
	"github.com/IBM/sarama"
	"github.com/go-keg/keg/contrib/config"
)

func NewSyncProducer(addrs []string, opts ...ConfigOption) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_1_1_0
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Producer.Partitioner = sarama.NewHashPartitioner
	for _, opt := range opts {
		opt(cfg)
	}

	client, err := sarama.NewSyncProducer(addrs, cfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewSyncProducerFromConfig(config config.Kafka, opts ...ConfigOption) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_1_1_0
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Producer.Partitioner = sarama.NewHashPartitioner

	opts = append(opts, SetNetSASL(config.User, config.Password))
	for _, opt := range opts {
		opt(cfg)
	}

	client, err := sarama.NewSyncProducer(config.GetAddr(), cfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}
