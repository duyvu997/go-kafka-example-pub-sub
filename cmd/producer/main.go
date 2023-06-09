package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/indralaksmana/queue_demo/model"
	kafka "github.com/indralaksmana/queue_demo/pkg/kafka"
	"github.com/sirupsen/logrus"
)

func main() {

	// Setup Logging
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)

	// connect to kafka
	kafkaBroker := []string{"127.0.0.1:9092"}
	kafkaProducer, errConnection := kafka.ConnectProducer(kafkaBroker)
	if errConnection != nil {
		logrus.Printf("error: %s", "Unable to configure kafka")
		return
	}
	defer kafkaProducer.Close()

	kafkaClientId := "1003"
	kafkaTopic := "test1"

	// send task to consumer via message broker
	message, errMarshal := json.Marshal(model.Message{
		Text: "Welcome to kafka in Golang",
	})

	if errMarshal != nil {
		logrus.Println(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": map[string]interface{}{
				"message": fmt.Sprintf("error while marshalling json: %s", errMarshal.Error()),
			},
		})
		return
	}

	errPushMessage := kafka.PushToQueue(kafkaBroker, kafkaClientId, kafkaTopic, message)
	if errPushMessage != nil {
		fmt.Println(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": map[string]interface{}{
				"message": fmt.Sprintf("error while push message into kafka: %s", errPushMessage.Error()),
			},
		})
		return
	}
}
