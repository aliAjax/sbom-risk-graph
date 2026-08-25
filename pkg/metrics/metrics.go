package metrics

import (
	"fmt"
	"sync/atomic"
)

type Counters struct {
	Imports, Evaluations, ParserErrors, Jobs, RecoveredPanics atomic.Uint64
}

func (c *Counters) Snapshot() map[string]uint64 {
	return map[string]uint64{"imports": c.Imports.Load(), "evaluations": c.Evaluations.Load(), "parser_errors": c.ParserErrors.Load(), "jobs": c.Jobs.Load(), "recovered_panics": c.RecoveredPanics.Load()}
}
func (c *Counters) RecordRecoveredPanic(committed bool) {
	if committed {
		return
	}
	c.RecoveredPanics.Add(1)
}
func Render(c *Counters) string {
	s := c.Snapshot()
	return fmt.Sprintf("sbom_imports_total %d\nsbom_evaluations_total %d\nsbom_parser_errors_total %d\nsbom_jobs_total %d\nsbom_recovered_panics_total %d\n", s["imports"], s["evaluations"], s["parser_errors"], s["jobs"], s["recovered_panics"])
}
