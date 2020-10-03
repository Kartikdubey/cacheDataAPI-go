package connection

import (
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

var KafkaConn *kafka.Conn
var err error

//KafkaConnection to establish connection
func KafkaConnection() {
	dialer := &kafka.Dialer{
		Timeout: 10 * time.Second,
	}
	for {
		KafkaConn, err = dialer.Dial("tcp", "your url")
		if err != nil {
			fmt.Println("Retring Kafka Connection")
			continue
		}
	}
}

//GetKafkaReader to get reader object for consuming messages
func GetKafkaReader(topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
}

//Similarly we can create writer objects,topics and use it for consuming and writing.
