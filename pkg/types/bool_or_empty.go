package types

import (
	"strconv"
	"strings"
)

type BoolOrEmpty struct {
	Valid bool
	Value bool
}

func (b *BoolOrEmpty) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "" {
		b.Valid = false
		b.Value = false
		return nil
	}

	v, err := strconv.ParseBool(str)
	if err != nil {
		return err
	}
	b.Valid = true
	b.Value = v
	return nil
}

func (b BoolOrEmpty) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte(`""`), nil
	}
	return []byte(strconv.FormatBool(b.Value)), nil
}
