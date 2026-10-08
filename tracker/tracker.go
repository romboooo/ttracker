package tracker

import "time"

type Session struct {
	Class     string
	StartTime time.Time
	EndTime   time.Time
}
