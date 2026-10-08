package newyeargolab1

import "time"

// DaysUntilNewYear возвращает количество календарных дней от даты t
// до 1 января следующего календарного года.
// Функция всегда смотрит вперёд: даже если t — это 1 января,
// результатом будет количество дней до Нового года следующего года.
// Учитываются високосные годы.
// Время суток (часы, минуты, секунды) игнорируется — учитываются
// только календарные дни. Расчёт ведётся в UTC, чтобы результат
// не зависел от переходов на летнее время.
func DaysUntilNextYear(t time.Time) int {
	t = t.UTC()
	startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	nextYear := time.Date(t.Year()+1, time.January, 1, 0, 0, 0, 0, time.UTC)
	return int(nextYear.Sub(startOfDay).Hours() / 24)
}
