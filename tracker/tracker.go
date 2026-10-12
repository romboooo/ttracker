package tracker

import (
	"fmt"

	"github.com/romboooo/ttracker/model"
)

type Session = model.Session

func PrintSession(session Session) {
	fmt.Printf("%s %v\n", "active window: "+session.Class+": ", session.EndTime.Sub(session.StartTime))
}
