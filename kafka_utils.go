package main

import (
	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
	"log"
	"router/pb"
)

const brokerAddress = "localhost:9092"

// Ideally set to > 1 for fault tolerance
const replicationFactor = 1

// Topic to publish analyses by Analyzer processes
const analysisStoreTopic = "topic.log.analysis.1"

// Topic to publish analysis requests, consumed by Analyzers
const analysisRequestStoreTopic = "topic.log.requests.analysis.1"

// Topic to publish init requests, consumed by Analyzers,
// to set up buffers and internal routing, single partition, because of less frequency
const initRequestStoreTopic = "topic.log.requests.init.1"

var brokerAddresses = []string{brokerAddress}

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
	createTopic(admin, initRequestStoreTopic, 1, replicationFactor, true)
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
	pkey := id2Service[msg.Id.Id]
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	message := &sarama.ProducerMessage{
		Topic: analysisRequestStoreTopic,
		Key:   sarama.StringEncoder(pkey),
		Value: sarama.ByteEncoder(data),
	}

	_, _, err = k.producer.SendMessage(message)
	return message, err
}

func (k *KafkaProducer) PublishInitRequest(msg *pb.InitRequest_Type0, _ map[string]string) (*sarama.ProducerMessage, error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	message := &sarama.ProducerMessage{
		Topic: initRequestStoreTopic,
		Value: sarama.ByteEncoder(data),
	}

	_, _, err = k.producer.SendMessage(message)
	return message, err
}
