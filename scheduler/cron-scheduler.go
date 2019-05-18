package scheduler

import (
	"fmt"
	"github.com/FrancescoIlario/Watering-Server/utils"
	"gopkg.in/robfig/cron.v2"
	"log"
)

var (
	c *cron.Cron = nil
)

func StartCron() {
	c = cron.New()
	_id, err := c.AddFunc("@every 10s", func() { fmt.Println("Every 10 seconds") })
	log.Printf("registered scheduler with entry %v\n", _id)
	_id2, err := c.AddFunc("@every 2m", func() { fmt.Println("Every 2 minutes") })
	log.Printf("registered scheduler with entry %v\n", _id2)
	utils.PanicIf(err)

	c.Start()
}


