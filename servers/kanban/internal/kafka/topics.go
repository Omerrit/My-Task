package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

func CompactTopicEntries() map[string]*string {
	var (
		cleanupPolicy    = "compact"
		deleteRetention  = "240000"
		minCompactionLag = "240000"
	)
	return map[string]*string{
		"delete.retention.ms":   &deleteRetention,
		"min.compaction.lag.ms": &minCompactionLag,
		"cleanup.policy":        &cleanupPolicy,
	}
}

func CheckTopic(topicName string, addresses []string, config *sarama.Config, topicEntries map[string]*string) error {
	client, err := sarama.NewClient(addresses, config)
	if err != nil {
		return err
	}

	topics, err := client.Topics()
	if err != nil {
		return err
	}
	brokers := client.Brokers()
	if len(brokers) == 0 {
		return fmt.Errorf("no brokers are avaliable")
	}
	broker := brokers[0]
	var topicExists bool
	for _, name := range topics {
		if name == topicName {
			topicExists = true
		}
	}
	if topicExists {
		request := &sarama.AlterConfigsRequest{
			Resources:    []*sarama.AlterConfigsResource{},
			ValidateOnly: false,
		}
		request.Resources = append(request.Resources, &sarama.AlterConfigsResource{
			Type:          sarama.TopicResource,
			Name:          topicName,
			ConfigEntries: topicEntries,
		})
		err := broker.Open(config)
		if err != nil {
			return err
		}
		resp, err := broker.AlterConfigs(request)
		if err != nil {
			return err
		}
		if len(resp.Resources) == 0 {
			return fmt.Errorf("unexpected response from kafka service")
		}
		if len(resp.Resources[0].ErrorMsg) > 0 {
			return fmt.Errorf("%s", resp.Resources[0].ErrorMsg)
		}
	} else {
		request := &sarama.CreateTopicsRequest{
			Version: 0,
			TopicDetails: map[string]*sarama.TopicDetail{
				topicName: {
					NumPartitions:     1,
					ReplicationFactor: -1,
					ReplicaAssignment: nil,
					ConfigEntries:     topicEntries,
				},
			},
			Timeout:      time.Second * 5,
			ValidateOnly: false,
		}
		err := broker.Open(config)
		if err != nil {
			return err
		}
		resp, err := broker.CreateTopics(request)
		if err != nil {
			return err
		}
		if len(resp.TopicErrors) == 0 {
			return fmt.Errorf("unexpected response from kafka service")
		}
		if resp.TopicErrors[topicName].ErrMsg != nil {
			return fmt.Errorf("%s", *resp.TopicErrors[topicName].ErrMsg)
		}
	}
	return nil
}
