package lab

import "time"

type Split struct {
	Name  string
	Start time.Time
	End   time.Time
}

func DefaultSplits() []Split {
	return []Split{
		{Name: "2021_2022_stress", Start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Name: "2023_2024_oos", Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Name: "2025_2026_recent", Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Name: "full_2021_2026"},
	}
}

func (s Split) Contains(t time.Time) bool {
	if s.Name == "full_2021_2026" {
		return true
	}
	return !t.Before(s.Start) && t.Before(s.End)
}
