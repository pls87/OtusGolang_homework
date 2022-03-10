package events

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
)

type EventIterator struct {
	rows *sqlx.Rows
}

func (s *EventIterator) Next() bool {
	return s.rows.Next()
}

func (s *EventIterator) Current() (models.Event, error) {
	var ev models.Event
	e := s.rows.StructScan(&ev)
	return ev, e
}

func (s *EventIterator) Complete() error {
	return s.rows.Close()
}

func (s *EventIterator) ToArray() ([]models.Event, error) {
	res := make([]models.Event, 0, 10)
	var ev models.Event
	for s.rows.Next() {
		e := s.rows.StructScan(&ev)
		if e != nil {
			_ = s.Complete()
			return nil, e
		}
	}

	return res, nil
}

// Execute TODO: clean up this code later.
func (ee *EventExpression) Execute(ctx context.Context) (basic.EventIterator, error) {
	clauseBuilder := make([]string, 0, 4)
	clauseArgs := make([]interface{}, 0, 9)
	ind := 0
	if ee.params.UserID > 0 {
		clauseBuilder = append(clauseBuilder, fmt.Sprintf("(user_id=$%d)", ind+1))
		clauseArgs = append(clauseArgs, ee.params.UserID)
		ind++
	}
	if !ee.params.ToNotify.IsZero() {
		clauseBuilder = append(clauseBuilder,
			fmt.Sprintf(`(start>$%d AND start - notify_before < $%d 
				AND "ID" NOT IN (SELECT event_id from "notification_sent"))`, ind+1, ind+2))
		clauseArgs = append(clauseArgs, ee.params.ToNotify, ee.params.ToNotify)
		ind += 2
	}
	if !ee.params.Starts.Start.IsZero() {
		clauseBuilder = append(clauseBuilder, fmt.Sprintf("(start>=$%d AND start<=$%d)", ind+1, ind+2))
		clauseArgs = append(clauseArgs, ee.params.Starts.Start, ee.params.Starts.End())
		ind += 2
	}

	if !ee.params.Intersection.Start.IsZero() {
		clauseBuilder = append(clauseBuilder,
			fmt.Sprintf("((start>=$%d AND start<=$%d) OR (start + duration >= $%d AND start + duration <= $%d))",
				ind+1, ind+2, ind+3, ind+4))
		clauseArgs = append(clauseArgs, ee.params.Intersection.Start, ee.params.Intersection.End(),
			ee.params.Intersection.Start, ee.params.Intersection.End())
	}

	whereClause := strings.Join(clauseBuilder, " AND ")
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}

	rows, err := ee.db.QueryxContext(ctx, `SELECT * FROM "events"`+whereClause, clauseArgs...) //nolint:sqlclosecheck
	if err != nil {
		return nil, err
	}

	return &EventIterator{rows}, nil
}
