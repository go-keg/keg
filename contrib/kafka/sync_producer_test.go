package kafka

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/IBM/sarama"
)

func TestNewSyncProducer(t *testing.T) {
	producer, err := NewSyncProducer(
		[]string{"127.0.0.1:9093"},
		SetNetSASL("", ""),
		SetProducerPartitioner(sarama.NewHashPartitioner),
	)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 30; i++ {
		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: "test",
			Key:   sarama.StringEncoder(strconv.Itoa(i)),
			Value: sarama.StringEncoder(fmt.Sprintf("test data %d", i)),
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("sync_producer: partition=%d offset=%d\n", partition, offset)
	}
}
