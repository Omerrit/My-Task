package kafka

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/errors"
	"github.com/Shopify/sarama"
	"log"
)

type consumer struct {
	actors.Actor
	consumer     sarama.Consumer
	partConsumer sarama.PartitionConsumer
	output       *consumerOutput
	name         string
	restored     bool
}

func NewConsumer(name string, system *actors.System, config *sarama.Config, brokers []string, topic string, partition int32, offset int64) (actors.ActorService, error) {
	service := new(consumer)
	var err error
	service.name = name
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, err
	}
	offsetNewest, err := client.GetOffset(topic, partition, sarama.OffsetNewest)
	if err != nil {
		return nil, err
	}
	offsetOldest, err := client.GetOffset(topic, partition, sarama.OffsetOldest)
	if err != nil {
		return nil, err
	}
	err = client.Close()
	if err != nil {
		return nil, err
	}
	service.restored = offsetNewest == offsetOldest
	service.consumer, err = sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}
	service.partConsumer, err = service.consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		return nil, err
	}
	return system.Spawn(service), nil
}

func (c *consumer) Run() error {
	for {
		select {
		case <-c.IncomingChannel():
			if !c.ProcessMessages() {
				return nil
			}
		case msg, ok := <-c.partConsumer.Messages():
			if !ok {
				return nil
			}
			c.output.Messages.Add(msg)
			if !c.restored && msg.Offset+1 == c.partConsumer.HighWaterMarkOffset() {
				c.output.Messages.Add(nil)
				c.restored = true
			}
			c.output.FlushLater()
			c.FlushReadyOutputs()
		}
	}
}

func (c *consumer) MakeBehaviour() actors.Behaviour {
	log.Printf("%s started", c.name)
	var behaviour actors.Behaviour
	c.SetPanicProcessor(c.onPanic)
	c.output = &consumerOutput{}
	c.output.CloseWhenActorCloses()
	behaviour.AddCommand(new(Subscribe), func(cmd interface{}) (actors.Response, error) {
		c.InitStreamOutput(c.output, cmd.(*Subscribe))
		return nil, nil
	})
	return behaviour
}

func (c *consumer) onPanic(err errors.StackTraceError) {
	log.Println("panic:", err, err.StackTrace())
	c.Quit(err)
}

func (c *consumer) Shutdown() error {
	err := c.partConsumer.Close()
	if err != nil {
		log.Println("error while shutting of partition consumer down:", err)
		return err
	}
	err = c.consumer.Close()
	if err != nil {
		log.Println("error while shutting of consumer down:", err)
		return err
	}
	log.Println(c.name, "shut down")
	return nil
}
