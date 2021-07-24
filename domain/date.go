package domain

import (
	"fmt"
	"time"
)

type YM struct {
	Year  int
	Month int
}

func NewYM(t time.Time) YM {
	return YM{
		Year:  t.Year(),
		Month: int(t.Month()),
	}
}

func (y *YM) String() string {
	return fmt.Sprintf("%d-%d", y.Year, y.Month)
}
