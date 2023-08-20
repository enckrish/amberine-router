package main

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
	"os"
	"router/pb"
)

// Ideally set to > 1 for fault tolerance
const replicationFactor = 1

// Topic to publish analyses by Analyzer processes
const analysisStoreTopic = "topic.log.analysis.result.1"

// Topic to publish analysis requests, consumed by Analyzers
const analysisRequestStoreTopic = "topic.log.requests.analysis.1"

var brokerAddresses = []string{os.Getenv("AMBER_KAFKA_URL")}

type KafkaProducer struct {
	config   *sarama.Config
	producer sarama.SyncProducer
}

func NewKafkaProducer(initTopic bool) *KafkaProducer {
	// Probably a better way to check topic's existence exists, if found, remove `initTopic`
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0

	// admin required for creating topics
	admin, err := sarama.NewClusterAdmin(brokerAddresses, config)
	defer func() { _ = admin.Close() }()
	if err != nil {
		log.Fatal("Error while creating cluster admin: ", err.Error())
	}
	if initTopic {
		createTopics(admin)
	}

	producer := createProducer(config)

	k := KafkaProducer{config: config, producer: producer}
	return &k
}

func createTopics(admin sarama.ClusterAdmin) {
	createTopic(admin, analysisStoreTopic, 3, replicationFactor, true)
	createTopic(admin, analysisRequestStoreTopic, 3, replicationFactor, true)
}

func createTopic(admin sarama.ClusterAdmin, topic string, numPartitions int32, replFactor int16, ignoreErr bool) {
	err := admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     numPartitions,
		ReplicationFactor: replFactor,
	}, false)
	if err != nil {
		if ignoreErr {
			return
		}
		log.Fatal("Error while creating topic: ", err.Error())
	}
}

func createProducer(config *sarama.Config) sarama.SyncProducer {
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokerAddresses, config)
	if err != nil {
		log.Fatal("Error while creating producer: ", err.Error())
	}
	return producer
}

func (k *KafkaProducer) PublishAnalysisRequest(msg *pb.AnalyzerRequest_Type0, id2Service map[string]string) (*sarama.ProducerMessage, error) {
	service := id2Service[msg.StreamId]
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	message := &sarama.ProducerMessage{
		Topic: analysisRequestStoreTopic,
		Key:   sarama.StringEncoder(service),
		Value: sarama.StringEncoder(data),
	}

	_, _, err = k.producer.SendMessage(message)
	return message, err
}

func (k *KafkaProducer) PublishInitRequest(msg *pb.InitRequest_Type0, _ map[string]string) (*sarama.ProducerMessage, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	message := &sarama.ProducerMessage{
		Topic: analysisRequestStoreTopic,
		Key:   sarama.StringEncoder(msg.Service),
		Value: sarama.StringEncoder(data),
	}

	_, _, err = k.producer.SendMessage(message)
	return message, err
}
