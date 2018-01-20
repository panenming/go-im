package every

import (
	"time"
)

type Expression struct {
	Interval time.Duration
}

func (e *Expression) Next(from time.Time) time.Time {
	return from.Add(e.Interval)
}

func NewExpression(dur time.Duration) *Expression {
	if dur < time.Second {
		dur = time.Second
	}
	return &Expression{
		Interval: dur,
	}
}
