package goutils

import (
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
