package canon

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func JSON(value any) ([]byte, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var normalized bytes.Buffer
	if err := json.Compact(&normalized, raw); err != nil {
		return nil, fmt.Errorf("canonical json: %w", err)
	}
	return normalized.Bytes(), nil
}
func Decode(data []byte, target any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(target); err != nil {
		return err
	}
	if dec.More() {
		return fmt.Errorf("trailing json")
	}
	return nil
}
