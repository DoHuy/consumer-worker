package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/Shopify/sarama"
	"io/ioutil"
	"log"
	"strings"
	"vtp/config"
)

type Kafka struct {
	Config		*config.KafkaConfig
	Client		sarama.Client
	Group		sarama.ConsumerGroup
	KafkaTopic	string
	KafkaChan	chan []byte
}
type consumerGroupHandler struct{ c chan []byte }

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.c <- msg.Value
		sess.MarkMessage(msg, "")
	}
	return nil
}

func New(config *config.KafkaConfig, topic, consumerGroup string) (*Kafka, error) {
	var kafkaInstance = &Kafka{Config: config}
	client, err := sarama.NewClient(strings.Split(config.KafkaBrokers, ","), kafkaInstance.getKafkaConfig())
	if err != nil {
		return nil, err
	}
	kafkaInstance.Client = client
	kafkaInstance.KafkaTopic = topic
	group, err := kafkaInstance.newConsumerGroupClient(consumerGroup)
	if err != nil {
		return nil, err
	}
	kafkaInstance.Group = group
	kafkaInstance.KafkaChan = make(chan []byte)
	return kafkaInstance, err
}

func (c *Kafka)newTLSConfig(clientCertFile, clientKeyFile, caCertFile string) (*tls.Config, error) {
	tlsConfig := tls.Config{}

	// Load client cert
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return &tlsConfig, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = caCertPool

	tlsConfig.BuildNameToCertificate()

	return &tlsConfig, err
}


func (c *Kafka)getKafkaConfig() *sarama.Config{
	configKafka := sarama.NewConfig()
	configKafka.Version = sarama.V1_0_0_0
	configKafka.Consumer.Return.Errors = true
	if c.Config.InitialOffset == "newest" {
		configKafka.Consumer.Offsets.Initial = sarama.OffsetNewest
	} else {
		configKafka.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	if c.Config.KafkaTLSEnabled {
		tlsConfig, err := c.newTLSConfig(c.Config.KafkaTLSClientCert, c.Config.KafkaTLSClientKey, c.Config.KafkaTLSCACertFile)
		if err != nil {
			log.Fatal(nil, "setup kafka TLS error", err)
		}
		tlsConfig.InsecureSkipVerify = true
		configKafka.Net.TLS.Enable = true
		configKafka.Net.TLS.Config = tlsConfig
	}
	return configKafka
}
func (c *Kafka) newConsumerGroupClient(kafkaConsumerGroup string) (sarama.ConsumerGroup, error){
	group, err := sarama.NewConsumerGroupFromClient(kafkaConsumerGroup, c.Client)
	if err != nil {
		log.Fatal(nil, "error on initialing kafka connection", err)
	}
	return group, err
}

func (c *Kafka)ConsumeGroup(ctx context.Context) error{
	handler := consumerGroupHandler{c.KafkaChan}
	topics := strings.Split(c.KafkaTopic, ",")
	err := c.Group.Consume(ctx, topics, handler)
	if err != nil {
		log.Println("ERROR =>", err)
	}
	return err
}
