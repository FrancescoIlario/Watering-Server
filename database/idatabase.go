package database

import (
	"github.com/FrancescoIlario/Watering-Server/schedule"
)

type IDbSchedule interface {
	Create(schedule schedule.Schedule)
	Read(id int) *schedule.Schedule
	ReadAll() []*schedule.Schedule
}

