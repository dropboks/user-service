package messagequeue

import (
	"context"
	"log"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type (
	NotificationStream interface {
		Publish(context.Context, string, []byte) (*jetstream.PubAck, error)
	}
	notificationStream struct {
		js     jetstream.JetStream
		logger zerolog.Logger
	}
)

func NewNotificationStream(js jetstream.JetStream, logger zerolog.Logger) (NotificationStream, error) {
	cfg := &jetstream.StreamConfig{
		Name:        viper.GetString("jetstream.notification.stream.name"),
		Description: viper.GetString("jetstream.notification.stream.description"),
		Subjects:    []string{viper.GetString("jetstream.notification.subject.global")},
		MaxBytes:    6 * 1024 * 1024,
		Storage:     jetstream.FileStorage,
	}
	_, err := js.CreateOrUpdateStream(context.Background(), *cfg)
	if err != nil {
		log.Printf("Failed to create or update JetStream Notification stream: %v", err)
		return nil, err
	}

	return &notificationStream{
		js:     js,
		logger: logger,
	}, nil
}

func (n *notificationStream) Publish(ctx context.Context, subject string, payload []byte) (*jetstream.PubAck, error) {
	ack, err := n.js.Publish(ctx, subject, payload)
	if err != nil {
		return nil, err
	}
	return ack, nil
}
