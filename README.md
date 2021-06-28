# How to run worker ?

Please export list env:

```export KAFKA_BROKERS="localhost:9093"
export KAFKA_TOPIC_TYPE1="sarama_topic_1"
export KAFKA_TOPIC_TYPE2="sarama_topic_2"
export KAFKA_CONSUMER_GROUP_TYPE_1="sarama_consumer"
export KAFKA_CONSUMER_GROUP_TYPE_2="sarama_consumer_2"
export AUTH_KAFKA_CONSUMER_INITIAL_OFFSET="newest"
export ENABLE_LOGGING="true"
export PRODUCE_ENDPOINT="http://api.vnsale.vn:8080/events/v1/invoice"
export ELASTIC_URLS="http://localhost:9200"
export VANDONHANHTRINH_INDEX="vandon_hanhtrinh"
export CHITIETDON_INDEX="chitiet_don"```