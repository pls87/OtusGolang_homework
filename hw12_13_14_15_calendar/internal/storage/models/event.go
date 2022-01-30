package models

import "time"

type ID uint64

const MaxDuration = time.Nanosecond * (1<<63 - 1)

type Timeframe struct {
	Start    time.Time     `db:"start" json:"start"`
	Duration time.Duration `db:"duration" json:"duration"`
}

func (t *Timeframe) End() time.Time {
	return t.Start.Add(t.Duration)
}

type Event struct {
	Timeframe
	ID           ID            `db:"ID" json:"id"`
	Title        string        `db:"title" json:"title"`
	UserID       ID            `db:"user_id" json:"user_id"`
	NotifyBefore time.Duration `db:"notify_before" json:"notify_before"`
	Desc         string        `db:"description" json:"desc"`
}
