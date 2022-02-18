package events

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
)

type EventExpression struct {
	params *basic.EventExpressionParams
	db     *sqlx.DB
}

func (ee *EventExpression) User(id models.ID) basic.EventExpression {
	ee.params.UserID = id
	return ee
}

func (ee *EventExpression) ToNotify() basic.EventExpression {
	ee.params.ToNotify = time.Now()
	return ee
}

func (ee *EventExpression) StartsIn(tf models.Timeframe) basic.EventExpression {
	ee.params.Starts = tf
	return ee
}

func (ee *EventExpression) Intersects(tf models.Timeframe) basic.EventExpression {
	ee.params.Intersection = tf
	return ee
}
