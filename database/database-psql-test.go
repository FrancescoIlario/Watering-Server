package database

import (
	"fmt"
	"github.com/FrancescoIlario/Watering-Server/schedule"
	"github.com/FrancescoIlario/Watering-Server/utils"
	"time"
)

func TestDb() {
	now := time.Now()
	year, month, day := now.Date()

	startOfToday := time.Date(year, month, day, 0, 0, 0, 0, now.Location())

	for i := 0; i < 5; i++ {
		slack, err := time.ParseDuration(fmt.Sprintf("%vm", i))
		utils.PanicIf(err)

		_d, err := time.ParseDuration("10m")
		utils.PanicIf(err)

		startDuration := now.Add(slack).Sub(startOfToday)
		endDuration := now.Add(slack).Add(_d).Sub(startOfToday)

		fmt.Printf("%v: startDuration %v, endDuration %v\n", i+1, startDuration, endDuration)

		newSchedule := schedule.Schedule{
			End:   time.Duration(endDuration),
			Start: time.Duration(startDuration),
		}

		utils.PanicIf(Create(&newSchedule))
	}

	fmt.Printf("\n****\n\n")

	schedules, err := ReadAll()
	utils.PanicIf(err)

	for _, sched := range schedules {
		fmt.Println(sched.String())
	}
	fmt.Printf("\nTotally read: %v\n", len(schedules))
}
