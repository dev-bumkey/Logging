package kafka

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cocktailcloud/acloud-alarm-collector/application/service"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/kafka"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/logger"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/types"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/utils"
	"github.com/zemirco/uid"
)

type Consumer struct {
	service service.AlarmService
	ready   chan bool
}

func New(conf *kafka.Config, svc service.AlarmService) {
	if !conf.UseKafka {
		logger.Info("Kafka is turned off")
		return
	}

	client, err := kafka.GetConsumerGroup(conf)
	if err != nil {
		panic(err)
	}
	// defer func() { _ = client.Close() }()
	consumer := Consumer{
		service: svc,
		ready:   make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, conf.Topic, &consumer); err != nil {
				logger.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	logger.Info("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		logger.Info("terminating: context cancelled")
	case <-sigterm:
		logger.Info("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}

}

func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(c.ready)
	return nil
}

func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	requestId := uid.New(10)

	for message := range claim.Messages() {
		now := time.Now()
		logger.Info("==========================================================================================")

		logger.Infof(">> message claimed at %s, reqeust-id: %s, topic: %s, message timestamp: %v",
			now.Format(utils.DefaultTimeFormat), requestId, message.Topic, message.Timestamp)
		reader, err := gzip.NewReader(bytes.NewReader(message.Value))
		if err != nil {
			logger.Error("fail to read metric from topic: ", err.Error())
			continue
		}

		buffer := bytes.Buffer{}
		read, err := buffer.ReadFrom(reader)
		if err != nil {
			logger.Error("fail to read alert from reader: ", err.Error())
			continue
		}
		logger.Debug("read bytes: ", read)

		// decoder := json.NewDecoder(bytes.NewReader(buffer.Bytes()))
		// var envelop types.TransmitEnvelop
		// if err = decoder.Decode(&envelop); err != nil {
		// 	logger.Error("fail to decode envelop: ", err.Error())
		// 	continue
		// }
		var transmit types.TransmitEnvelop
		errs := json.Unmarshal(buffer.Bytes(), &transmit)
		if errs != nil {
			logger.Errorf("transmit Unmarshal Error >> %v", errs)
			// return errs
			continue
		}

		logger.Info(fmt.Sprintf("   - Agent Info >> Cluster: %v | Transmit: %v | TransmitterVersion: %v | Audit Message Size: %v ",
			transmit.Cluster, transmit.Transmitter, transmit.TransmitterVersion, len(transmit.Audits)))

		// if proc.collectorConfig.DevMode && proc.collectorConfig.SkipSave {
		// 	logger.Debug("skip save metrics")
		// 	continue
		// }

		//Database 저장 처리
		// err = execute.InsertMetrics(c.service, transmit.Cluster, &transmit)
		// _, err = c.service.InsertAlerts(&transmit)
		_, err = c.service.ReceiveAlarmsProcess(&transmit)
		if err != nil {
			logger.Errorf("   - Database Info >> InsertAlerts Error Message : %v", err)
			// return errs
			continue
		}

		session.MarkMessage(message, "")
		logger.Infof(">> message claimed processed request-id: %s, elapsed time: %v\n\n", requestId, time.Since(now))
	}
	logger.Debug("consumer ConsumeClaim end")

	return nil
}
