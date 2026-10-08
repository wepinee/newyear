package newyear

import (
	"testing"
	"time"
)

func TestDaysUntilNewYear(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want int
	}{
		{
			name: "начало года — 1 января 2025",
			date: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			want: 365, // 2025 — невисокосный, значит от 1 янв 2025 до 1 янв 2026 = 365
		},
		{
			name: "конец года — 31 декабря 2025",
			date: time.Date(2025, time.December, 31, 0, 0, 0, 0, time.UTC),
			want: 1,
		},
		{
			name: "високосный год — 1 января 2024",
			date: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			want: 366, // 2024 — високосный
		},
		{
			name: "день до 29 февраля високосного года",
			date: time.Date(2024, time.February, 28, 0, 0, 0, 0, time.UTC),
			want: 308, // от 28 фев 2024 до 1 янв 2025
		},
		{
			name: "29 февраля високосного года",
			date: time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC),
			want: 307,
		},
		{
			name: "день после 29 февраля високосного года",
			date: time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC),
			want: 306,
		},
		{
			name: "игнорирование времени суток — 31 декабря 2025, 23:59",
			date: time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC),
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DaysUntilNextYear(tt.date)
			if got != tt.want {
				t.Errorf("DaysUntilNewYear(%v) = %d, want %d", tt.date, got, tt.want)
			}
		})
	}
}
