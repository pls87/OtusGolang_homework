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

func (t *Timeframe) Period(d time.Time, str string) (ok bool) {
	switch str {
	case "day":
		t.Day(d)
	case "week":
		t.Week(d)
	case "month":
		t.Month(d)
	default:
		return false
	}
	return true
}

func (t *Timeframe) Day(d time.Time) {
	t.Start = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	t.Duration = 24 * time.Hour
}

func (t *Timeframe) Week(d time.Time) {
	dd := d
	for dd.Weekday() != time.Monday {
		dd = dd.AddDate(0, 0, -1)
	}
	t.Start = time.Date(dd.Year(), dd.Month(), dd.Day(), 0, 0, 0, 0, time.UTC)
	t.Duration = 7 * 24 * time.Hour
}

func (t *Timeframe) Month(d time.Time) {
	t.Start = time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, time.UTC)
	t.Start.AddDate(0, 1, 0)
	t.Duration = t.Start.AddDate(0, 1, 0).Sub(t.Start)
}

type Event struct {
	Timeframe
	ID           ID            `db:"ID" json:"id"`
	Title        string        `db:"title" json:"title"`
	UserID       ID            `db:"user_id" json:"user_id"`             //nolint:tagliatelle
	NotifyBefore time.Duration `db:"notify_before" json:"notify_before"` //nolint:tagliatelle
	Desc         string        `db:"description" json:"desc"`
}
