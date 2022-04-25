package catalog

import (
	"fmt"
	"time"
)

type metrics struct {
	Find, Done, Error     int
	StartDate, FinishDate time.Time
}

func (m *metrics) IncreaseFind()  { m.Find += 1 }
func (m *metrics) IncreaseDone()  { m.Done += 1 }
func (m *metrics) IncreaseError() { m.Error += 1 }

func (m *metrics) GetStatistic() string {
	temp := `
	Найдено: %d
	Обработано: %d
	Ошибки: %d
	Время выполнения: %v`
	return fmt.Sprintf(temp, m.Find, m.Done, m.Error, m.GetWastedTime())
}

func (m *metrics) GetWastedTime() time.Duration {
	return m.FinishDate.Sub(m.StartDate)
}
