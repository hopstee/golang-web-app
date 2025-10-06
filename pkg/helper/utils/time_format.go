package utils

import "time"

type TimeFormat string

const (
	BaseFormat         TimeFormat = "02.01.2006"
	BaseFormatWithTime TimeFormat = "02.01.2006 15:04"
)

func FormatTimeString(ts string, format TimeFormat) string {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return ts
	}

	if format == "" {
		format = BaseFormat
	}

	return t.Format(string(format))
}
