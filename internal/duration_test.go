package internal

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name string
		d    time.Duration
		want string
	}{
		{"100ms", 100 * time.Millisecond, "100ms"},
		{"1s", time.Second, "1s"},
		{"minute", time.Minute, "60s"},
		{"90s", 90 * time.Second, "90s"},
		{"2h", 2 * time.Hour, "7200s"},
		{"1500ms", 1500 * time.Millisecond, "1500ms"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDuration(tt.d); got != tt.want {
				t.Errorf("FormatDuration(%v) = %q, want %q", tt.d, got, tt.want)
			}
		})
	}
}
