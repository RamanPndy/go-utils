package goutils_test

import (
	"testing"
	"time"

	goutils "github.com/RamanPndy/go-utils/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestUnixTimeToTimestamp(t *testing.T) {
	tests := []struct {
		name     string
		unixTime int64
		wantTime time.Time
	}{
		{
			name:     "Standard timestamp",
			unixTime: 1634235600, // Equivalent to 2021-10-14T00:00:00Z
			wantTime: time.Unix(1634235600, 0),
		},
		{
			name:     "Unix epoch",
			unixTime: 0, // Equivalent to 1970-01-01T00:00:00Z
			wantTime: time.Unix(0, 0),
		},
		{
			name:     "Far future timestamp",
			unixTime: 32503680000, // Equivalent to 3000-01-01T00:00:00Z
			wantTime: time.Unix(32503680000, 0),
		},
		{
			name:     "Far past timestamp",
			unixTime: -2208988800, // Equivalent to 1900-01-01T00:00:00Z
			wantTime: time.Unix(-2208988800, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := goutils.UnixTimeToTimestamp(tt.unixTime)

			// Verify if the converted timestamp matches the expected result
			want := timestamppb.New(tt.wantTime)
			if !got.AsTime().Equal(want.AsTime()) {
				t.Errorf("UnixTimeToTimestamp(%d) = %v, want %v", tt.unixTime, got.AsTime(), want.AsTime())
			}
		})
	}
}
