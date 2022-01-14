package models

import "time"

type ID uint64

const MaxDuration = time.Nanosecond * (1<<63 - 1)

type Timeframe struct {
	Start    time.Time
	Duration time.Duration
}

type Event struct {
	Timeframe
	ID           ID
	Title        string
	UserID       ID
	NotifyBefore time.Duration
	Desc         string
}
