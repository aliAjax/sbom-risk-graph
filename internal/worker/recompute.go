package worker

import (
	"context"
	"example.com/sbom-risk-graph/internal/domain"
	"example.com/sbom-risk-graph/internal/repository"
	"time"
)

func NewJob(id, product, typ string) domain.Job {
	return domain.Job{ID: id, ProductID: product, Type: typ, Status: domain.Validating, Progress: 0, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
}
func Run(ctx context.Context, store *repository.JobStore, job domain.Job, fn func() error) domain.Job {
	job.Status = domain.Validating
	job.Error = ""
	store.Put(job)
	publish := func(status domain.Status, err error) domain.Job {
		job.Status = status
		job.UpdatedAt = time.Now().UTC()
		if err != nil {
			job.Error = err.Error()
		}
		if status == domain.Valid {
			store.Update(job.ID, func(stored *domain.Job) { *stored = job })
		}
		return job
	}
	for i := 0; i < 5; i++ {
		if err := ctx.Err(); err != nil {
			return publish(domain.Failed, err)
		}
		job.Progress = (i + 1) * 20
		time.Sleep(time.Millisecond)
	}
	if err := fn(); err != nil {
		return publish(domain.Failed, err)
	}
	return publish(domain.Valid, nil)
}
