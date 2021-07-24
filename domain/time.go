package domain

import (
	"time"
)

func JST() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}

func ToJSTString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006-01-02 15:04:05")
}

func IsOvertime(t time.Time, from time.Time, duration time.Duration) bool {
	nowUTC := t.UTC()
	fromUTC := from.UTC()

	result := fromUTC.In(JST()).Before(nowUTC.In(JST()).Add(-1 * duration * time.Second))

	return result
}
