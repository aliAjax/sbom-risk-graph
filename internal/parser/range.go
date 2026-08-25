package parser

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidRange = errors.New("invalid affected range")

type Comparator struct {
	Op      string
	Version Version
}
type Range struct{ Comparators []Comparator }

func ParseRange(raw string) (Range, error) {
	var r Range
	for _, part := range strings.Fields(raw) {
		op := "="
		value := part
		if strings.HasPrefix(part, ">=") {
			op, value = ">=", part[2:]
		} else if strings.HasPrefix(part, "<=") {
			op, value = "<=", part[2:]
		} else if strings.HasPrefix(part, ">") || strings.HasPrefix(part, "<") || strings.HasPrefix(part, "=") {
			op, value = part[:1], part[1:]
		}
		v, err := ParseVersion(value)
		if err != nil {
			return Range{}, err
		}
		r.Comparators = append(r.Comparators, Comparator{Op: op, Version: v})
	}
	if len(r.Comparators) == 0 {
		return Range{}, fmt.Errorf("empty range")
	}
	return r, nil
}
func (r Range) Matches(v Version) bool {
	for _, c := range r.Comparators {
		cmp := v.Compare(c.Version)
		ok := map[string]bool{"=": cmp == 0, ">": cmp > 0, "<": cmp < 0, ">=": cmp >= 0, "<=": cmp <= 0}[c.Op]
		if !ok {
			return false
		}
	}
	return true
}
