package newyear

import (
	"encoding/json"
	"net/http"
	"time"
)

// dateFormat — формат даты, который API принимает и возвращает.
// Соответствует ISO 8601 (date-only).
const dateFormat = time.DateOnly

// DaysResponse — тело успешного ответа GET /days.
type DaysResponse struct {
	// Date — дата, от которой производился расчёт, в формате YYYY-MM-DD.
	Date string `json:"date"`
	// Days — количество календарных дней до 1 января следующего года.
	Days int `json:"days"`
}

// ErrorResponse — тело ответа при ошибке.
type ErrorResponse struct {
	// Error — человекочитаемое описание ошибки.
	Error string `json:"error"`
}

// DaysHandler возвращает http.Handler, обслуживающий GET /days.
//
// Поведение:
//   - без параметра date — расчёт ведётся относительно текущей даты (UTC);
//   - с параметром date=YYYY-MM-DD — расчёт ведётся относительно указанной даты;
//   - при некорректной дате возвращается 400 Bad Request;
//   - при неверном HTTP-методе возвращается 405 Method Not Allowed.
func DaysHandler() http.Handler {
	return http.HandlerFunc(handleDays)
}

func handleDays(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed, use GET")
		return
	}

	dateStr := r.URL.Query().Get("date")

	var date time.Time
	if dateStr == "" {
		date = time.Now().UTC()
	} else {
		parsed, err := time.Parse(dateFormat, dateStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
			return
		}
		date = parsed
	}

	resp := DaysResponse{
		Date: date.Format(dateFormat),
		Days: DaysUntilNextYear(date),
	}
	writeJSON(w, http.StatusOK, resp)
}

// writeJSON сериализует v в JSON и записывает его в w с указанным статусом.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError отправляет JSON-ответ с описанием ошибки.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}
