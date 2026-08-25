package parser

import "strings"

func DetectFormat(data []byte) string {
	value := strings.ToLower(strings.TrimSpace(string(data)))
	if strings.Contains(value, "cyclonedx") || strings.Contains(value, "bomformat") {
		return "cyclonedx"
	}
	if strings.Contains(value, "spdxversion") || strings.Contains(value, "spdxid") {
		return "spdx"
	}
	return "unknown"
}
