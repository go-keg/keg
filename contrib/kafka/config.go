package kafka

import "github.com/IBM/sarama"

type ConfigOption func(config *sarama.Config)

func SetProducerPartitioner(partitioner sarama.PartitionerConstructor) ConfigOption {
	return func(config *sarama.Config) {
		config.Producer.Partitioner = partitioner
	}
}

func SetConsumerOffsetInitial(offset int64) ConfigOption {
	return func(config *sarama.Config) {
		config.Consumer.Offsets.Initial = offset
	}
}

func SetConsumerGroupBalanceStrategy(strategy ...sarama.BalanceStrategy) ConfigOption {
	return func(config *sarama.Config) {
		config.Consumer.Group.Rebalance.GroupStrategies = strategy
	}
}

func SetNetSASL(user, password string) ConfigOption {
	return func(config *sarama.Config) {
		if user != "" && password != "" {
			config.Net.SASL.Enable = true
			config.Net.SASL.User = user
			config.Net.SASL.Password = password
		}
	}
}

func SetVersion(version sarama.KafkaVersion) ConfigOption {
	return func(config *sarama.Config) {
		config.Version = version
	}
}

func SetConsumerFetchMax(n int32) ConfigOption {
	return func(config *sarama.Config) {
		config.Consumer.Fetch.Max = n
	}
}
