package attestation

import "context"

type Attestation struct {
	ID        string
	ProductID string
	Digest    string
}

type Source interface {
	Next(context.Context) (Attestation, error)
}

func readNext(ctx context.Context, source Source) (Attestation, error) {
	return source.Next(context.Background())
}
