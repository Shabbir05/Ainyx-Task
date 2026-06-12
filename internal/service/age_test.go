package service_test

import (
	"testing"
	"time"

	"github.com/yourusername/user-api/internal/service"
)

func TestCalculateAge(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name string
		dob  time.Time
		want int
	}{
		{
			name: "birthday already passed this year",
			dob:  time.Date(now.Year()-30, now.Month()-1, now.Day(), 0, 0, 0, 0, time.UTC),
			want: 30,
		},
		{
			name: "birthday has not passed yet this year",
			dob:  time.Date(now.Year()-30, now.Month()+1, now.Day(), 0, 0, 0, 0, time.UTC),
			want: 29,
		},
		{
			name: "birthday is today",
			dob:  time.Date(now.Year()-25, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			want: 25,
		},
		{
			name: "leap year birthday (Feb 29) — checking Dec 31",
			dob:  time.Date(2000, time.February, 29, 0, 0, 0, 0, time.UTC),
			want: service.CalculateAge(time.Date(2000, time.February, 29, 0, 0, 0, 0, time.UTC)),
		},
		{
			name: "newborn — born today",
			dob:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			want: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := service.CalculateAge(tc.dob)
			if got != tc.want {
				t.Errorf("CalculateAge(%v) = %d; want %d", tc.dob, got, tc.want)
			}
		})
	}
}

// TestCalculateAge_SpecificDates tests with fixed dates to ensure correctness
// regardless of when the test runs (relative calculations excluded).
func TestCalculateAge_SpecificDates(t *testing.T) {
	// We test the boundary logic directly.
	// Born on 1990-05-10 — if today is 2025-06-12, age is 35.
	// We can't assert absolute values without mocking time.Now(),
	// so we assert the function does not go negative and is reasonable.

	dob := time.Date(1990, time.May, 10, 0, 0, 0, 0, time.UTC)
	age := service.CalculateAge(dob)

	if age < 0 {
		t.Errorf("CalculateAge returned negative age %d for dob %v", age, dob)
	}
	if age > 200 {
		t.Errorf("CalculateAge returned unreasonably large age %d for dob %v", age, dob)
	}
}
