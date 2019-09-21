package gcal

import (
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
)

type EventsInWeek map[time.Weekday][]*calendar.Event

type Calendar struct {
	ID  string
	Svc *calendar.Service
}

// day=n: n日後の予定を取得
func (c *Calendar) ListEventsInNDays(n int) ([]*calendar.Event, error) {
	nextDay := time.Now().AddDate(0, 0, n)

	s := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, time.Local).Format(time.RFC3339)
	e := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 23, 59, 59, 999999999, time.Local).Format(time.RFC3339)

	events, err := c.Svc.Events.List(c.ID).ShowDeleted(false).
		SingleEvents(true).TimeMin(s).TimeMax(e).OrderBy("startTime").Do()
	if err != nil {
		return nil, err
	}

	return events.Items, err
}

func (c *Calendar) BatchGetEventsInNextWeek() (EventsInWeek, error) {
	now := time.Now()
	var dltDay int
	if now.Weekday() == time.Sunday {
		dltDay = 1
	} else {
		dltDay = 8 - int(now.Weekday())
	}

	// 時間をセットするために一時変数が必要
	mtmp := now.AddDate(0, 0, dltDay)
	stmp := mtmp.AddDate(0, 0, 6)

	monday := time.Date(mtmp.Year(), mtmp.Month(), mtmp.Day(), 0, 0, 0, 0, time.Local).Format(time.RFC3339)
	sunday := time.Date(stmp.Year(), stmp.Month(), stmp.Day(), 23, 59, 59, 999999999, time.Local).Format(time.RFC3339)

	events, err := c.Svc.Events.List(c.ID).ShowDeleted(false).
		SingleEvents(true).TimeMin(monday).TimeMax(sunday).OrderBy("startTime").Do()
	if err != nil {
		return nil, err
	}

	ret := make(EventsInWeek)
	for _, i := range events.Items {
		date := i.Start.Date
		if date == "" {
			date = strings.Split("T", i.Start.DateTime)[0]
		}
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, err
		}
		ret[t.Weekday()] = append(ret[t.Weekday()], i)
	}
	return ret, nil
}
