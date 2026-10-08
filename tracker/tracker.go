package tracker

import (
	"fmt"
	"time"
)

type Session struct {
	Class     string
	StartTime time.Time
	EndTime   time.Time
}

func PrintSession(session Session) {
	fmt.Printf("%s %v\n", "active window: "+session.Class+": ", session.EndTime.Sub(session.StartTime))
}
