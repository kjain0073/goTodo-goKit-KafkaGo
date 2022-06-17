package view

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	logO "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/kjain0073/go-Todo/models"
	"github.com/kjain0073/go-Todo/tasks"

	"github.com/segmentio/kafka-go"
	"gopkg.in/mgo.v2/bson"
)

const (
	topic         = "tasks-topic1"
	brokerAddress = "localhost:9092"
)

func SaveToKafka(ctx context.Context, todo models.TodoEntity) error {
	fmt.Println("saving to kafka")
	todoDto := models.TodoDto{
		ID:        todo.ID.Hex(),
		Title:     todo.Title,
		Completed: todo.Completed,
		CreatedAt: todo.CreatedAt,
	}
	jsonString, err := json.Marshal(todoDto)
	if err != nil {
		panic(err)
	}

	taskString := string(jsonString)

	l := log.New(os.Stdout, "kafka writer: ", 0)
	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		// assign the logger to the writer
		Logger: l,
	})

	for _, word := range []string{string(taskString)} {
		// each kafka message has a key and value. The key is used
		// to decide which partition (and consequently, which broker)
		// the message gets published on
		err := w.WriteMessages(ctx, kafka.Message{
			Key: []byte(todoDto.ID),
			// create an arbitrary message payload for the value
			Value: []byte(word),
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}

		// log a confirmation once the message is written
		fmt.Println("writes:", taskString)
		// sleep for a second
		time.Sleep(time.Second)
	}
	return nil
}

func PollFromKafka(s tasks.Service, ctx context.Context, logger logO.Logger) error {
	fmt.Println("Start receiving from Kafka")

	// c.Close()
	l := log.New(os.Stdout, "kafka reader: ", 0)
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		GroupID: "my-group",
		// assign the logger to the reader
		Logger: l,
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value

		taskString := string(msg.Value)
		var task models.TodoDto
		b := []byte(taskString)
		err2 := json.Unmarshal(b, &task)
		fmt.Println("value is " + taskString)

		if err2 != nil {
			fmt.Println(err2)
		}
		// to poll into db
		var todo = models.TodoEntity{
			ID:        bson.ObjectIdHex(task.ID),
			Title:     task.Title,
			Completed: task.Completed,
			CreatedAt: task.CreatedAt,
		}

		fmt.Println("received: ", todo.Title)
		if _, err3 := s.CreateTodoToRepo(ctx, todo); err3 != nil {
			level.Error(logger).Log("err", err3)
			return err3
		}
	}

}
