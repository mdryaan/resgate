package priority

import "github.com/mdryaan/resgate/internal/config"

type Level int

const (
	Min Level = config.MinPriority
	Max Level = config.MaxPriority
)

func (p Level) IsHigherThan(other Level) bool { return p > other }
func (p Level) IsLowerThan(other Level) bool  { return p < other }

func (p Level) IsValid() bool {
	return int(p) >= config.MinPriority && int(p) <= config.MaxPriority
}
