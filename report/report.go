package report

import (
	"context"
	"sort"
	"time"

	"github.com/romboooo/ttracker/model"
)

type SessionStore interface {
	SessionsBetween(ctx context.Context, from, to time.Time) ([]model.Session, error)
}
type SessionDuration struct {
	Class    string
	Duration time.Duration
}

func Today(ctx context.Context, store SessionStore) ([]SessionDuration, error) {
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfToday := startOfToday.AddDate(0, 0, 1)
	sessions, err := store.SessionsBetween(ctx,
		startOfToday,
		endOfToday,
	)

	if err != nil {
		return nil, err
	}
	return SortByDuration(SummarizeByClass(sessions, startOfToday, endOfToday)), nil
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

func SummarizeByClass(sessions []model.Session, from, to time.Time) map[string]time.Duration {
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

func DurationInRange(session model.Session, from, to time.Time) time.Duration {
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
