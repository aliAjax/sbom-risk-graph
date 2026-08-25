package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major, Minor, Patch int
	Raw                 string
}

func ParseVersion(raw string) (Version, error) {
	clean := strings.TrimPrefix(strings.TrimSpace(raw), "v")
	parts := strings.Split(clean, ".")
	if len(parts) < 1 || len(parts) > 3 {
		return Version{}, fmt.Errorf("invalid semantic version %q", raw)
	}
	v := Version{Raw: raw}
	var values [3]int
	for i, p := range parts {
		n, err := strconv.Atoi(strings.FieldsFunc(p, func(r rune) bool { return r < '0' || r > '9' })[0])
		if err != nil {
			return Version{}, fmt.Errorf("invalid version %q", raw)
		}
		values[i] = n
	}
	v.Major, v.Minor, v.Patch = values[0], values[1], values[2]
	return v, nil
}
func (v Version) Compare(other Version) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}
	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}
	if v.Patch < other.Patch {
		return -1
	}
	if v.Patch > other.Patch {
		return 1
	}
	return 0
}
func (v Version) String() string { return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch) }
