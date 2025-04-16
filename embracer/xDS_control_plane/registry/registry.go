package registry

import (
	"embracer/utils"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Registry struct {
	ingestC chan *kgo.Record
}

func NewRegistry(ingest chan *kgo.Record) Registry {
	return Registry{ingestC: ingest}
}

func (s Registry) StartRegistry() error {

	registry := utils.Registry{
		Id:      "1",
		Status:  "start",
		Message: "xDS_control_plane start!",
	}

	rb, err := json.Marshal(registry)
	if err != nil {
		log.Error().Msgf("%v", err)
		return err
	}

	s.ingestC <- &kgo.Record{
		Topic:   "StartRegistry",
		Value:   rb,
		Headers: []kgo.RecordHeader{{Key: "namespace", Value: []byte("1")}},
	}
	log.Info().Msgf("[Registry] StartRegistry -> topic[StartRegistry]")

	return nil
}
