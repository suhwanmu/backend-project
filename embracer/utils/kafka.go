package utils

import (
	"context"
	"embracer/utils/log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

// kafka client logger
var (
	KgoLogger kgo.Logger = kgo.BasicLogger(log.NewIOWriter(zerolog.InfoLevel), kgo.LogLevelInfo, nil)
)

type Registry struct {
	Id      string `json:"id,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

func CheckKafkaConn(kafkaBrokers []string) error {
	opts := []kgo.Opt{
		kgo.SeedBrokers(kafkaBrokers...),
		kgo.WithLogger(KgoLogger),
	}
	kClient, err := kgo.NewClient(opts...)
	if err != nil {
		return err
	}
	defer kClient.Close()
	if err := kClient.Ping(context.Background()); err != nil {
		return err
	}
	return nil
}

func StartKafka(
	cfg Config,
) (chan *kgo.Record, context.CancelFunc) {

	topic := []string{"StartRegistry"}

	if err := CreateTopics(topic, cfg); err != nil {
		log.Error().Msg(err.Error())
	}
	kafkaCtx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM)

	// this for sending messages to kafka
	IngestC := make(chan *kgo.Record, 10000)
	go StartProducer(kafkaCtx, cfg.KafkaBrokers, IngestC)
	return IngestC, cancel
}

func StartProducer(
	ctx context.Context,
	brokers []string,
	ingestChan chan *kgo.Record,
) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.WithLogger(KgoLogger),
		kgo.UnknownTopicRetries(3),
		kgo.RecordRetries(10),
		kgo.AutoCommitInterval(1 * time.Second),
	}

	kClient, err := kgo.NewClient(opts...)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer kClient.Close()

	if err := kClient.Ping(ctx); err != nil {
		log.Error().Msg(err.Error())
	}

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("stop producing to kafka")
			if err := kClient.Flush(context.Background()); err != nil {
				log.Error().Msg(err.Error())
			}
			return

		case record := <-ingestChan:
			kClient.Produce(
				ctx,
				record,
				func(r *kgo.Record, err error) {
					if err != nil {
						log.Error().Msgf(
							"failed to produce record %s record: %v",
							err, record,
						)
					}
				},
			)
		}
	}

}

func CreateTopics(topics []string, cfg Config) error {
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.KafkaBrokers...),
		kgo.WithLogger(KgoLogger),
	}

	kClient, err := kgo.NewClient(opts...)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer kClient.Close()

	if err := kClient.Ping(context.Background()); err != nil {
		log.Error().Msg(err.Error())
	}

	adminClient := kadm.NewClient(kClient)
	defer adminClient.Close()

	topicConfig := map[string]*string{
		"retention.ms": kadm.StringPtr(cfg.KafkaTopicRetentionMs),
	}
	resp, err := adminClient.CreateTopics(context.Background(),
		cfg.KafkaTopicPartitions, cfg.KafkaTopicReplicas, topicConfig, topics...)
	if err != nil {
		return err
	}
	for _, r := range resp.Sorted() {
		if r.Err != nil {
			log.Error().Msgf("topic: %s error: %s", r.Topic, r.Err)
		}
	}
	return nil
}
