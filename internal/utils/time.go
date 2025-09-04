package utils

import (
	"fmt"
	"time"
)

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}

func TimeUntil(t *time.Time) string {
	if t == nil {
		return "never"
	}
	remaining := time.Until(*t)
	if remaining <= 0 {
		return "expired"
	}
	return FormatDuration(remaining)
}
