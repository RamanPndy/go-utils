package goutils

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UnixTimeToTimestamp(unixTime int64) *timestamp.Timestamp {
	// Convert unix time to time.Time
	ut := time.Unix(unixTime, 0)

	// Convert time.Time to Timestamp
	ts := timestamppb.New(ut)

	return ts
}

func NormalizeTime(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return t.UTC().Round(0)
}

func ParseDate(raw any) (time.Time, error) {
	if t, ok := raw.(time.Time); ok {
		return t.UTC(), nil
	}
	text, ok := raw.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("invalid date")
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return time.Time{}, fmt.Errorf("empty date")
	}
	if t, err := time.Parse("2006-01-02", text); err == nil {
		return t.UTC(), nil
	}
	t, err := time.Parse(time.RFC3339, text)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

func ParseDateTime(raw any) (time.Time, error) {
	if t, ok := raw.(time.Time); ok {
		return t.UTC(), nil
	}
	text, ok := raw.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("invalid datetime")
	}
	t, err := time.Parse(time.RFC3339, strings.TrimSpace(text))
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}
