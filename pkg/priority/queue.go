package priority

import (
	"sort"

	"github.com/mdryaan/resgate/internal/models"
)

type Entry struct {
	Reservation *models.Reservation
	TenantName  string
	Priority    Level
}

type Queue struct {
	entries []Entry
}

func New() *Queue {
	return &Queue{}
}

func (q *Queue) Push(r *models.Reservation, p Level) {
	q.entries = append(q.entries, Entry{Reservation: r, TenantName: r.TenantName, Priority: p})
	sort.Slice(q.entries, func(i, j int) bool {
		return q.entries[i].Priority > q.entries[j].Priority
	})
}

func (q *Queue) Pop() (Entry, bool) {
	if len(q.entries) == 0 {
		return Entry{}, false
	}
	e := q.entries[0]
	q.entries = q.entries[1:]
	return e, true
}

func (q *Queue) Peek() (Entry, bool) {
	if len(q.entries) == 0 {
		return Entry{}, false
	}
	return q.entries[0], true
}

func (q *Queue) Len() int { return len(q.entries) }

func (q *Queue) RemoveTenant(tenantName string) []Entry {
	var removed []Entry
	remaining := q.entries[:0]
	for _, e := range q.entries {
		if e.TenantName == tenantName {
			removed = append(removed, e)
		} else {
			remaining = append(remaining, e)
		}
	}
	q.entries = remaining
	return removed
}
