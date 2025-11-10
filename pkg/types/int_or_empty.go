package types

import (
	"strconv"
	"strings"
)

type IntOrEmpty struct {
	Valid bool
	Value int
}

// Для JSON
func (i *IntOrEmpty) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "" {
		i.Valid = false
		i.Value = 0
		return nil
	}

	num, err := strconv.Atoi(str)
	if err != nil {
		return err
	}
	i.Valid = true
	i.Value = num
	return nil
}

// Для маршала, если нужно
func (i IntOrEmpty) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte(`""`), nil
	}
	return []byte(strconv.Itoa(i.Value)), nil
}
