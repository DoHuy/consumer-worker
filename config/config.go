package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	KafkaConfig   *KafkaConfig   `envconfig:"KAFKA_CONFIG" required:"true"`
	VnSaleAPIs    *VnSaleAPIs    `envconfig:"VNSALE_APIS" required:"true"`
	ElasticConfig *ElasticConfig `envconfig:"ELASTIC_CONFIG" required:"true"`
	IsProduction  bool           `envconfig:"IS_PRODUCTION" default:"false"`
}
type ElasticConfig struct {
	URLs                 string `envconfig:"ELASTIC_URLS" required:"true"`
	IndexVanDonHanhTrinh string `envconfig:"VANDONHANHTRINH_INDEX" required:"true"`
	IndexChiTietDon      string `envconfig:"CHITIETDON_INDEX" required:"true"`
}
type KafkaConfig struct {
	KafkaBrokers            string `envconfig:"KAFKA_BROKERS" required:"true"`
	KafkaTopicType1         string `envconfig:"KAFKA_TOPIC_TYPE1" required:"true"`
	KafkaTopicType2         string `envconfig:"KAFKA_TOPIC_TYPE2" required:"true"`
	InitialOffset           string `envconfig:"AUTH_KAFKA_CONSUMER_INITIAL_OFFSET" default:"newest"`
	KafkaConsumerGroupType1 string `envconfig:"KAFKA_CONSUMER_GROUP_TYPE_1" required:"true"`
	KafkaConsumerGroupType2 string `envconfig:"KAFKA_CONSUMER_GROUP_TYPE_2" required:"true"`
	KafkaTLSCACertFile      string `envconfig:"KAFKA_TLS_CA_CERT_FILE" default:"/secrets/ca_cert.pem"`
	KafkaTLSClientCert      string `envconfig:"KAFKA_TLS_CLIENT_CERT" default:"/secrets/client_cert.pem"`
	KafkaTLSClientKey       string `envconfig:"KAFKA_TLS_CLIENT_KEY" default:"/secrets/client_key.pem"`
	KafkaTLSEnabled         bool   `envconfig:"KAFKA_TLS_ENABLED" default:"false"`
}

type VnSaleAPIs struct {
	ProduceEndpoint string `envconfig:"PRODUCE_ENDPOINT" default:"http://api.vnsale.vn:8080/events/v1/invoice"`
}

func GetBasicConfig() *Config {
	workerConfig := Config{}
	if err := envconfig.Process("", &workerConfig); err != nil {
		log.Fatal(nil, "load env error", err)
	}
	return &workerConfig

}
