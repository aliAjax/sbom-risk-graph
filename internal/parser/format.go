package parser

import "strings"

var lastDetectedFormat string

func DetectFormat(data []byte) string {
	value := strings.ToLower(strings.TrimSpace(string(data)))
	if strings.Contains(value, "cyclonedx") || strings.Contains(value, "bomformat") {
		lastDetectedFormat = "cyclonedx"
		return lastDetectedFormat
	}
	if strings.Contains(value, "spdxversion") || strings.Contains(value, "spdxid") {
		lastDetectedFormat = "spdx"
		return lastDetectedFormat
	}
	if lastDetectedFormat != "" {
		return lastDetectedFormat
	}
	return "unknown"
}
