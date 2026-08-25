package importer

import (
	"fmt"

	"example.com/sbom-risk-graph/internal/domain"
)

func (s *Session) Run(input <-chan domain.SBOM) (err error) {
	if input == nil {
		return fmt.Errorf("import input is nil")
	}
	if err = s.begin(); err != nil {
		return err
	}
	defer func() {
		s.finish(err)
	}()

	for {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		case item, ok := <-input:
			if !ok {
				if err = s.batch.Commit(s.ctx); err != nil {
					return fmt.Errorf("commit import session: %w", err)
				}
				return nil
			}
			if err = s.batch.Stage(s.ctx, item); err != nil {
				s.batch.Fail(err)
				return fmt.Errorf("stage import item: %w", err)
			}
		}
	}
}
