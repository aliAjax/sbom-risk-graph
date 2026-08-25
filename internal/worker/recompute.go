package worker

import (
	"context"
	"example.com/sbom-risk-graph/internal/domain"
	"time"
)

func NewJob(id, product, typ string) domain.Job {
	return domain.Job{ID: id, ProductID: product, Type: typ, Status: domain.Validating, Progress: 0, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
}
func Run(ctx context.Context, job domain.Job, fn func() error) domain.Job {
	job.Status = domain.Validating
	for i := 0; i < 5; i++ {
		if err := ctx.Err(); err != nil {
			job.Status = domain.Failed
			job.Error = err.Error()
			return job
		}
		job.Progress = (i + 1) * 20
		time.Sleep(time.Millisecond)
	}
	if err := fn(); err != nil {
		job.Status = domain.Failed
		job.Error = err.Error()
		return job
	}
	job.Status = domain.Valid
	job.UpdatedAt = time.Now().UTC()
	return job
}
