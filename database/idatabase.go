package database

import (
	"github.com/filario/watering-server/schedule"
)

type IDbSchedule interface {
	Create(schedule schedule.Schedule)
	Read(id int) *schedule.Schedule
	ReadAll() []*schedule.Schedule
}

