package conversion

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseBool(v interface{}) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case string:
		lower := strings.ToLower(strings.TrimSpace(val))
		return lower == "true" || lower == "1" || lower == "yes"
	case float64:
		return val != 0
	case int:
		return val != 0
	case int64:
		return val != 0
	default:
		return false
	}
}

func ParseString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case fmt.Stringer:
		return val.String()
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", val)
	}
}

func ParseInt(v interface{}) (int, bool) {
	if v == nil {
		return 0, false
	}
	switch val := v.(type) {
	case int:
		return val, true
	case int64:
		return int(val), true
	case float64:
		return int(val), true
	case string:
		val = strings.TrimSpace(val)
		if val == "" {
			return 0, false
		}
		f, err := strconv.Atoi(val)
		return f, err == nil
	case bool:
		if val {
			return 1, true
		}
		return 0, true
	default:
		return 0, false
	}
}

func ParseFloat(v interface{}) (float64, bool) {
	if v == nil {
		return 0, false
	}
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case string:
		val = strings.TrimSpace(val)
		if val == "" {
			return 0, false
		}
		f, err := strconv.ParseFloat(val, 64)
		return f, err == nil
	case bool:
		if val {
			return 1, true
		}
		return 0, true
	default:
		return 0, false
	}
}
