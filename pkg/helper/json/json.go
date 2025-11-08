package json

import (
	"encoding/json"
	"fmt"
)

func MapToStruct(src interface{}, dest interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err := json.Unmarshal(bytes, dest); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}
