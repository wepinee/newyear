package newyear

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDaysHandler(t *testing.T) {
	srv := httptest.NewServer(DaysHandler())
	defer srv.Close()

	client := srv.Client()

	tests := []struct {
		name       string
		query      string
		wantStatus int
		wantDays   int
		wantDate   string
		wantError  bool
	}{
		{
			name:       "конец года — 31 декабря 2025",
			query:      "?date=2025-12-31",
			wantStatus: http.StatusOK,
			wantDays:   1,
			wantDate:   "2025-12-31",
		},
		{
			name:       "начало года — 1 января 2025",
			query:      "?date=2025-01-01",
			wantStatus: http.StatusOK,
			wantDays:   365,
			wantDate:   "2025-01-01",
		},
		{
			name:       "високосный год — 1 января 2024",
			query:      "?date=2024-01-01",
			wantStatus: http.StatusOK,
			wantDays:   366,
			wantDate:   "2024-01-01",
		},
		{
			name:       "без параметра date — расчёт от текущей даты",
			query:      "",
			wantStatus: http.StatusOK,
		},
		{
			name:       "некорректный формат даты",
			query:      "?date=abc",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name:       "некорректный формат — 2025/12/31",
			query:      "?date=2025/12/31",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Get(srv.URL + tt.query)
			if err != nil {
				t.Fatalf("GET failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			if tt.wantError {
				var e ErrorResponse
				if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
					t.Fatalf("decode error response: %v", err)
				}
				if e.Error == "" {
					t.Errorf("expected non-empty error message")
				}
				return
			}

			var got DaysResponse
			if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
				t.Fatalf("decode response: %v", err)
			}

			if tt.wantDate != "" && got.Date != tt.wantDate {
				t.Errorf("date = %q, want %q", got.Date, tt.wantDate)
			}
			if tt.wantDays != 0 && got.Days != tt.wantDays {
				t.Errorf("days = %d, want %d", got.Days, tt.wantDays)
			}
			if tt.wantDate == "" {
				// Кейс без параметра: дата должна быть сегодняшней (UTC).
				today := time.Now().UTC().Format(dateFormat)
				if got.Date != today {
					t.Errorf("date = %q, want today %q", got.Date, today)
				}
			}
		})
	}
}

func TestDaysHandlerWrongMethod(t *testing.T) {
	srv := httptest.NewServer(DaysHandler())
	defer srv.Close()

	resp, err := srv.Client().Post(srv.URL, "application/json", nil)
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
	}
	if allow := resp.Header.Get("Allow"); allow != http.MethodGet {
		t.Errorf("Allow = %q, want %q", allow, http.MethodGet)
	}
}
