package infra_event

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"
	event_broker "weather_notification/src/modules/notification_schedule/domain/event"
	infra_database "weather_notification/src/modules/shared/infra/database"

	"github.com/segmentio/kafka-go"
)

type KafkaEventBroker struct {
	database  infra_database.Database
	kafkaConn *kafka.Conn
	handlers  map[string][]func(message []byte) error
	topic     string
	context   context.Context
	host      []string
	partition int
}

func (broker KafkaEventBroker) Emit(event event_broker.Event) error {
	broker.kafkaConn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	eventPayload := reflect.ValueOf(event.Payload)

	switch eventPayload.Kind() {
	case reflect.Array, reflect.Slice:
		messages := make([]kafka.Message, 0, eventPayload.Len())

		for i := 0; i < eventPayload.Len(); i++ {
			payload := eventPayload.Index(i).Interface()

			jsonData, err := json.Marshal(payload)

			if err != nil {
				return errors.New("could not parse data")
			}

			messages = append(messages, kafka.Message{Value: jsonData})
		}
		broker.kafkaConn.WriteMessages(messages...)
	case reflect.Struct:
		jsonData, err := json.Marshal(event.Payload)

		if err != nil {
			return errors.New("could not parse data")
		}

		broker.kafkaConn.WriteMessages(kafka.Message{Value: jsonData})
	default:
		return fmt.Errorf("unsupported type: %s", eventPayload.Kind().String())
	}

	return nil
}

func (broker KafkaEventBroker) Subscribe(eventName string, handler func(message []byte) error) error {
	broker.handlers[eventName] = append(broker.handlers[eventName], handler)
	return nil
}

func (broker KafkaEventBroker) Listen() {
	go func(ctx context.Context) {
		broker.kafkaConn.SetDeadline(time.Now().Add(10 * time.Second))

		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:   broker.host,
			Topic:     broker.topic,
			Partition: broker.partition,
			MaxBytes:  10e6,
		})

		topicTable := "topic_partition_offset"

		lastPartitionOffsetQuery := broker.database.QueryBuilder(topicTable).
			SetColumns([]string{"offset"}).Where("topic = ?").Select()

		var offset int64

		topicName := fmt.Sprintf("%s-%d", broker.topic, broker.partition)

		row := broker.database.SelectOne(lastPartitionOffsetQuery, topicName)

		if err := row.Scan(&offset); err != nil {
			if err != sql.ErrNoRows {
				log.Fatal("could not read last offset", err)
				return
			}

			insertQuery := broker.database.QueryBuilder(topicTable).
				SetColumns([]string{"topic", "offset"}).Insert(1)

			if err := broker.database.Exec(insertQuery, topicName, 0); err != nil {
				log.Fatal("could not initialize partition offset", err)
			}

			offset = 0
		}

		reader.SetOffset(offset)

		for {
			m, err := reader.ReadMessage(ctx)

			if err != nil {
				break
			}

			for _, handler := range broker.handlers[broker.topic] {
				if err := handler(m.Value); err != nil {
					fmt.Println(err)
					continue
				}
			}

			updateQuery := broker.database.QueryBuilder(topicTable).
				SetColumns([]string{"offset"}).
				Where("topic = ?").
				Update()

			offset += 1

			if err := broker.database.Exec(updateQuery, offset, topicName); err != nil {
				log.Fatal("could not save partition", err)
			}
		}

		if err := reader.Close(); err != nil {
			return
		} else {
			log.Fatal("could not close reader")
		}
	}(broker.context)

}

func NewKafkaEventBroker(topic string, partition int, host string, context context.Context, database infra_database.Database) KafkaEventBroker {
	conn, err := kafka.DialLeader(context, "tcp", host, topic, partition)

	if err != nil {
		log.Fatal("failed to connect to kafka", err)
	}

	return KafkaEventBroker{
		kafkaConn: conn,
		handlers:  make(map[string][]func(message []byte) error),
		topic:     topic,
		context:   context,
		partition: partition,
		host:      []string{host},
		database:  database,
	}
}
