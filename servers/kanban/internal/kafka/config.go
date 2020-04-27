package kafka

import "github.com/Shopify/sarama"

func NewConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_4_0_0
	return config
}
