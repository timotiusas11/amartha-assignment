package nsq

import (
	"encoding/json"
	"fmt"
)

// "github.com/nsqio/go-nsq"

type NSQInterface interface {
	Send(channel string, value interface{}) error
}

type NSQ struct {
	// Client
}

func New() NSQ {
	// config := nsq.NewConfig()
	// w, err := nsq.NewProducer("127.0.0.1:4150", config)
	// if err != nil {
	// 	log.Panic("Error while connecting producer to nsqd")
	// }
	return NSQ{
		// Set client
	}
}

func (n NSQ) Send(channel string, value interface{}) error {
	bvalue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	// n.Client.Send(...)
	fmt.Println("Message sent!")
	fmt.Println(string(bvalue))
	return nil
}
