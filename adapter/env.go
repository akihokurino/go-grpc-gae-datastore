package adapter

import (
	"time"
)

type SwitchProvider struct {
	IsMaintenance bool
}

type ThresholdProvider struct {
	NoMessageDuration time.Duration
	NoEntryDuration   time.Duration
}
