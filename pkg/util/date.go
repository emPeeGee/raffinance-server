package util

import "time"

func EndOfTheDay(date time.Time) *time.Time {
	endOfTheDay := date.
		Truncate(24 * time.Hour).
		Add(time.Hour*23 + time.Minute*59 + time.Second*59 + time.Millisecond*999)

	return &endOfTheDay
}
