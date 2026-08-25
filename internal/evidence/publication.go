package evidence

import (
	"context"

	"example.com/sbom-risk-graph/internal/domain"
)

type transitionRegistry interface {
	Transition(string, domain.Status) error
}

type Publisher struct {
	registry transitionRegistry
}

func NewPublisher(registry transitionRegistry) *Publisher {
	return &Publisher{registry: registry}
}

func (p *Publisher) Publish(ctx context.Context, id string, next domain.Status) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return p.registry.Transition(id, next)
}
