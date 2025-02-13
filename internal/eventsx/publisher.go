package eventsx

import (
	"context"

	"go.infratographer.com/x/events"
)

// Publisher represents something that sends relationships to permissions-api via NATS.
type Publisher interface {
	// PublishAuthRelationship is similar to events.Publisher.PublishAuthRelationship, but with no topic.
	PublishAuthRelationshipRequest(ctx context.Context, message events.AuthRelationshipRequest) (events.Message[events.AuthRelationshipResponse], error)
}

type publisher struct {
	topic    string
	innerPub events.AuthRelationshipPublisher
}

func (p publisher) PublishAuthRelationshipRequest(ctx context.Context, message events.AuthRelationshipRequest) (events.Message[events.AuthRelationshipResponse], error) {
	return p.innerPub.PublishAuthRelationshipRequest(ctx, p.topic, message)
}

// NewPublisher creates a new events publisher from the given config.
func NewPublisher(cfg Config) (Publisher, error) {
	natsCfg := events.NATSConfig{
		URL:           cfg.NATS.URL,
		PublishPrefix: cfg.NATS.PublishPrefix,
		Token:         cfg.NATS.Token,
		CredsFile:     cfg.NATS.CredsFile,
	}

	conn, err := events.NewNATSConnection(natsCfg)
	if err != nil {
		return nil, err
	}

	out := publisher{
		topic:    cfg.NATS.PublishTopic,
		innerPub: conn,
	}

	return out, nil
}
