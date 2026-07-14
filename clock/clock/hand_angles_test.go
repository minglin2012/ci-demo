package clock

import (
	"math"
	"testing"
	"time"
)

func TestCalculateAngles(t *testing.T) {
	tests := []struct {
		name                       string
		timeStr                    string
		wantHour, wantMin, wantSec float64
	}{
		{
			name:    "12:00:00 — all hands at 12",
			timeStr: "2026-07-15T12:00:00Z",
			wantHour: 0, wantMin: 0, wantSec: 0,
		},
		{
			name:    "12:15:00 — quarter past",
			timeStr: "2026-07-15T12:15:00Z",
			wantHour: 7.5, wantMin: 90, wantSec: 0,
		},
		{
			name:    "12:00:30 — 30 seconds",
			timeStr: "2026-07-15T12:00:30Z",
			wantHour: 0, wantMin: 3, wantSec: 180,
		},
		{
			name:    "6:00:00 — 6 o'clock",
			timeStr: "2026-07-15T06:00:00Z",
			wantHour: 180, wantMin: 0, wantSec: 0,
		},
		{
			name:    "3:30:00 — half past three",
			timeStr: "2026-07-15T03:30:00Z",
			wantHour: 105, wantMin: 180, wantSec: 0,
		},
		{
			name:    "00:00:00 — midnight",
			timeStr: "2026-07-15T00:00:00Z",
			wantHour: 0, wantMin: 0, wantSec: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := time.Parse(time.RFC3339, tt.timeStr)
			if err != nil {
				t.Fatalf("failed to parse time: %v", err)
			}

			h, m, s := CalculateAngles(parsed)

			if !almostEqual(h, tt.wantHour) {
				t.Errorf("hour angle = %.1f, want %.1f", h, tt.wantHour)
			}
			if !almostEqual(m, tt.wantMin) {
				t.Errorf("minute angle = %.1f, want %.1f", m, tt.wantMin)
			}
			if !almostEqual(s, tt.wantSec) {
				t.Errorf("second angle = %.1f, want %.1f", s, tt.wantSec)
			}
		})
	}
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.001
}

func TestCalculateAnglesContinuousSecond(t *testing.T) {
	base, _ := time.Parse(time.RFC3339, "2026-07-15T12:00:00Z")
	for sec := 0; sec < 60; sec++ {
		tm := base.Add(time.Duration(sec) * time.Second)
		_, m, s := CalculateAngles(tm)

		expectedSec := float64(sec) * 6.0
		expectedMin := float64(sec) * 0.1

		if !almostEqual(s, expectedSec) {
			t.Errorf("second %d: second angle = %.1f, want %.1f", sec, s, expectedSec)
		}
		if !almostEqual(m, expectedMin) {
			t.Errorf("second %d: minute angle = %.1f, want %.1f", sec, m, expectedMin)
		}
	}
}
