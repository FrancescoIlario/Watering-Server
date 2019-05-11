package schedule

import (
	"fmt"
	"time"
)

type Schedule struct {
	Id int64

	End   time.Duration
	Start time.Duration
}

func (schedule *Schedule) InTime(duration time.Duration) bool {
	return schedule.Start > duration && schedule.End < duration
}

func (schedule Schedule) String() string {
	return fmt.Sprintf("Id: %v, Start: %v, End: %v", schedule.Id, schedule.Start, schedule.End)
}
