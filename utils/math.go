package goutils

import (
	"encoding/json"
	"math"
	"strings"
)

func AsFloat(v any) (float64, bool) {
	switch value := v.(type) {
	case float64:
		return value, true
	case float32:
		return float64(value), true
	case int:
		return float64(value), true
	case int8:
		return float64(value), true
	case int16:
		return float64(value), true
	case int32:
		return float64(value), true
	case int64:
		return float64(value), true
	case uint:
		return float64(value), true
	case uint8:
		return float64(value), true
	case uint16:
		return float64(value), true
	case uint32:
		return float64(value), true
	case uint64:
		return float64(value), true
	case json.Number:
		f, err := value.Float64()
		if err != nil {
			return 0, false
		}
		return f, true
	case string:
		var n json.Number = json.Number(strings.TrimSpace(value))
		f, err := n.Float64()
		if err != nil {
			return 0, false
		}
		return f, true
	default:
		return 0, false
	}
}

func Round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func Round4(v float64) float64 {
	return math.Round(v*10000) / 10000
}
