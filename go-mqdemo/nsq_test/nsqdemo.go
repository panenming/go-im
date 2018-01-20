package main

import (
	"fmt"

	"github.com/nsqio/go-nsq"
)

const msgCount = 10

// 消费者
type ConsumerHandler struct {
	name string
}

func (consumer *ConsumerHandler) HandleMessage(msg *nsq.Message) error {
	fmt.Println(consumer.name, string(msg.Body))
	return nil
}

// 生产者
func Producer() {
	producer, err := nsq.NewProducer("127.0.0.1:4150", nsq.NewConfig())
	if err != nil {
		fmt.Println("NewProducer", err)
		panic(err)
	}
	i := 1
	for {
		if err := producer.Publish("test", []byte(fmt.Sprintf("hello world %d", i))); err != nil {
			fmt.Println("Publish", err)
			panic(err)
		}
		//time.Sleep(5 * time.Second)
		i++
		if i == msgCount {
			return
		}
	}
}

// 消费者A
func ConsumerA() {
	consumer, err := nsq.NewConsumer("test", "test-channel-a", nsq.NewConfig())
	if err != nil {
		fmt.Println("NewConsumerA", err)
		panic(err)
	}

	consumer.AddHandler(&ConsumerHandler{
		name: "ConsumerA",
	})

	if err := consumer.ConnectToNSQLookupd("127.0.0.1:4161"); err != nil {
		fmt.Println("ConnectToNSQLookupd", err)
		panic(err)
	}
}

// ConsumerB 消费者
func ConsumerB() {
	consumer, err := nsq.NewConsumer("test", "test-channel-b", nsq.NewConfig())
	if err != nil {
		fmt.Println("NewConsumerB", err)
		panic(err)
	}

	consumer.AddHandler(&ConsumerHandler{
		name: "ConsumerB",
	})

	if err := consumer.ConnectToNSQLookupd("127.0.0.1:4161"); err != nil {
		fmt.Println("ConnectToNSQLookupd", err)
		panic(err)
	}
}

func main() {
	//ConsumerA()
	ConsumerB()
	//Producer()
	for {
	}
}
