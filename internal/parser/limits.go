package parser

import "fmt"

type Limits struct {
	MaxComponents   int
	MaxDependencies int
	MaxDepth        int
}

func DefaultLimits() Limits {
	return Limits{MaxComponents: 100000, MaxDependencies: 300000, MaxDepth: 64}
}
func (l Limits) Validate(componentCount, dependencyCount int) error {
	if componentCount > l.MaxComponents {
		return fmt.Errorf("component count exceeds %d", l.MaxComponents)
	}
	if dependencyCount > l.MaxDependencies {
		return fmt.Errorf("dependency count exceeds %d", l.MaxDependencies)
	}
	return nil
}
