package main

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"log"
	"strconv"
	"time"
	"vtp/config"
	"vtp/dao"
	"vtp/dto"
	"vtp/es"
	"vtp/kafka"
	"vtp/logger"
	"vtp/utils"
)

func main() {
	workerConfig := config.GetBasicConfig()
	//init zap log
	var logger, err = logger.NewLoggerService(workerConfig.IsProduction)
	if err != nil {
		panic(err)
	}
	// init elastic search
	var esClient, errEs = es.New(workerConfig.ElasticConfig.URLs, workerConfig.ElasticConfig.Index)
	if errEs != nil {
		panic(errEs)
	}
	// init instance kafka
	kafkaInstanceType1, err := kafka.New(workerConfig.KafkaConfig, workerConfig.KafkaConfig.KafkaTopicType1, workerConfig.KafkaConfig.KafkaConsumerGroupType1)
	kafkaInstanceType2, err := kafka.New(workerConfig.KafkaConfig, workerConfig.KafkaConfig.KafkaTopicType2, workerConfig.KafkaConfig.KafkaConsumerGroupType2)
	if err != nil {
		panic(err)
	}
	//routine consume topic type 1
	ctx := context.Background()
	go func() {
		eventChan := make(chan dto.VandonhanhtrinhType1)
		// init 2 goroutine send events
		go func() {
			for d := range kafkaInstanceType1.KafkaChan {
				var eventData dto.VandonhanhtrinhType1
				err := json.Unmarshal(d, &eventData)
				if err != nil {
					log.Println("unmarshal error", err)
					func() { _ = kafkaInstanceType1.Group.Close() }()
					func() { _ = kafkaInstanceType1.Client.Close() }()
				} else {
					eventChan <- eventData
				}
			}
		}()
		go func() {
			for e := range eventChan {
				// todo: logic implement to do here
				var document dao.Vandonhanhtrinh
				if err := utils.JsonToJson(e.Data, &document); err != nil {
					logger.Error("json to json failed ", zap.Error(err))
					continue
				}
				rawBody, _ := json.Marshal(document)
				if err := esClient.Insert(ctx, string(rawBody)); err != nil {
					logger.Error("insert to elk failed", zap.Error(err))
					continue
				}
				logger.Info("insert success to es", zap.String("data", string(rawBody)))
				// push to VN sale
				var invoice dto.Invoice
				invoice.ID = document.OrderNumber
				invoice.Time = time.Unix(document.TrackingTime, 0)
				i, _ := strconv.ParseInt(fmt.Sprintf("%v", document.PostId), 10, 64)
				invoice.PostID = i
				invoice.ShipperID = fmt.Sprintf("%v", document.Employee)
				invoice.State = document.OrderStatus
				invoice.Description = fmt.Sprintf("%v", document.ProductDescription)
				invoiceBody, _ := json.Marshal(invoice)
				_, err := utils.MakePOSTRequestAPI(workerConfig.VnSaleAPIs.ProduceEndpoint, string(invoiceBody))
				if err != nil {
					logger.Error("request to VN sale failed", zap.Error(err))
					continue
				}
				logger.Info("push to VN sale success", zap.String("data", string(invoiceBody)))
			}
		}()

		logger.Info("start listening van don hanh trinh events type 1 ...")
		for {
			err := kafkaInstanceType1.ConsumeGroup(ctx)
			if err != nil {
				func() { _ = kafkaInstanceType1.Group.Close() }()
				func() { _ = kafkaInstanceType1.Client.Close() }()
				log.Fatal(err)

			}
		}
	}()

	// main consume topic type 2
	contextMain := context.Background()
	var eventChanType2 = make(chan dto.VandonhanhtrinhType2)
	// init 2 goroutine send events
	go func() {
		for d := range kafkaInstanceType2.KafkaChan {
			var eventData dto.VandonhanhtrinhType2
			fmt.Println("event data => ", string(d))
			err := json.Unmarshal(d, &eventData)
			if err != nil {
				log.Println("unmarshal error", err)
				func() { _ = kafkaInstanceType2.Group.Close() }()
				func() { _ = kafkaInstanceType2.Client.Close() }()

			} else {
				eventChanType2 <- eventData
			}
		}
	}()
	go func() {
		for e := range eventChanType2 {
			// todo: logic implement to do here
			var document dao.Vandonhanhtrinh
			if err := utils.JsonToJson(e.Data, &document); err != nil {
				logger.Error("json to json failed ", zap.Error(err))
				continue
			}
			rawBody, _ := json.Marshal(document)
			if err := esClient.Insert(contextMain, string(rawBody)); err != nil {
				logger.Error("insert to elk failed", zap.Error(err))
				continue
			}
			logger.Info("insert success to es", zap.String("data", string(rawBody)))
			// push to VN sale
			var invoice dto.Invoice
			invoice.ID = document.OrderNumber
			invoice.Time = time.Unix(document.TrackingTime, 0)
			i, _ := strconv.ParseInt(fmt.Sprintf("%v", document.PostId), 10, 64)
			invoice.PostID = i
			invoice.ShipperID = fmt.Sprintf("%v", document.Employee)
			invoice.State = document.OrderStatus
			invoice.Description = fmt.Sprintf("%v", document.ProductDescription)
			invoiceBody, _ := json.Marshal(invoice)
			_, err := utils.MakePOSTRequestAPI(workerConfig.VnSaleAPIs.ProduceEndpoint, string(invoiceBody))
			if err != nil {
				logger.Error("request to VN sale failed", zap.Error(err))
				continue
			}
			logger.Info("push to VN sale success", zap.String("data", string(invoiceBody)))
		}
	}()

	logger.Info("start listening van don hanh trinh events type 2 ...")
	for {
		err := kafkaInstanceType2.ConsumeGroup(contextMain)
		if err != nil {
			func() { _ = kafkaInstanceType2.Group.Close() }()
			func() { _ = kafkaInstanceType2.Client.Close() }()
			log.Fatal(err)

		}
	}
}

