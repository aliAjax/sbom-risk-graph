package attestation

import (
	"context"
	"fmt"
)

type Pipeline struct {
	ctx       context.Context
	cancel    context.CancelFunc
	collector *Collector
}

func NewPipeline(parent context.Context) *Pipeline {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	return &Pipeline{ctx: ctx, cancel: cancel, collector: NewCollector()}
}

func (p *Pipeline) Close()                { p.cancel() }
func (p *Pipeline) Done() <-chan struct{} { return p.ctx.Done() }

func (p *Pipeline) Import(source Source) (Receipt, error) {
	item, err := readNext(p.ctx, source)
	if err != nil {
		p.collector.Abort()
		return Receipt{State: "failed"}, wrapImportError(err)
	}
	if item.ID == "" || item.ProductID == "" {
		p.collector.Abort()
		return Receipt{State: "failed"}, wrapImportError(fmt.Errorf("attestation id and product id are required"))
	}
	p.collector.Stage(item)
	if err := p.collector.Commit(); err != nil {
		return Receipt{State: "failed"}, wrapImportError(err)
	}
	return Receipt{State: "committed", IDs: p.collector.IDs()}, nil
}
