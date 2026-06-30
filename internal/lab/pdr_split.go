package lab

import "time"

func btc15MPrevDaySplit(t time.Time, splits []Split) string {
	for _, split := range splits {
		if split.Contains(t) {
			return split.Name
		}
	}
	return "unassigned"
}
