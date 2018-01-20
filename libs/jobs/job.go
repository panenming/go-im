package jobs

import (
	"sort"
	"time"
)

type job struct {
	Function   func()
	Expression Expression
	Name       string
	Previous   time.Time
	Next       time.Time
}

type jobs []*job

func (j jobs) Len() int {
	return len(j)
}

func (j jobs) Swap(x, y int) {
	j[x], j[y] = j[y], j[x]
}

func (j jobs) Less(x, y int) bool {
	if j[x].Next.IsZero() {
		return false
	}
	if j[y].Next.IsZero() {
		return true
	}
	return j[x].Next.Before(j[y].Next)
}

func (j jobs) Sort() {
	sort.Sort(j)
}
