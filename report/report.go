package report

import (
	"sort"
	"time"

	"github.com/romboooo/ttracker/tracker"
)

type SessionDuration struct {
	Class    string
	Duration time.Duration
}

func SortByDuration(summary map[string]time.Duration) []SessionDuration {

	result := make([]SessionDuration, 0, len(summary))

	for class, dur := range summary {
		result = append(result, SessionDuration{
			Class:    class,
			Duration: dur,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Duration > result[j].Duration
	})

	return result

}

func SummarizeByClass(sessions []tracker.Session, from, to time.Time) map[string]time.Duration {
	sessionToTime := make(map[string]time.Duration)

	for _, session := range sessions {
		sessionToTime[session.Class] += DurationInRange(session, from, to)
	}

	for class, dur := range sessionToTime {
		if dur == 0 {
			delete(sessionToTime, class)
		}
	}
	return sessionToTime
}

func DurationInRange(session tracker.Session, from, to time.Time) time.Duration {
	startTime := session.StartTime
	endTime := session.EndTime

	if from.After(startTime) {
		startTime = from
	}

	if to.Before(endTime) {
		endTime = to
	}

	if !endTime.After(startTime) {
		return 0
	}

	return endTime.Sub(startTime)
}
