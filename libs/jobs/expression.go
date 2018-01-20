package jobs

import (
	"time"
)

type Expression interface {
	Next(time.Time) time.Time
}
