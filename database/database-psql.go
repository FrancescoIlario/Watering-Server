package database

import (
	. "github.com/filario/watering-server/schedule"
	. "github.com/filario/watering-server/utils"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func newDb() *pg.DB {
	return pg.Connect(&pg.Options{
		User: "postgres",
		Password: "",
		Addr: "127.0.0.1:5432",
		Database: "watering",
	})
}

func Initialize() {
	_db := newDb()
	defer _db.Close()

	for _, model := range []interface{}{(*Schedule)(nil)} {
		err := _db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:true,
		})
		PanicIf(err)
	}
}

func Create(schedule *Schedule) error {
	conn := newDb()
	defer conn.Close()

	return conn.Insert(&schedule)
}

func Read(id int64) (*Schedule, error) {
	conn := newDb()
	defer conn.Close()

	readSchedule := new(Schedule)
	err := conn.
		Model(readSchedule).
		Where("id = ?", id).
		Select()

	return readSchedule, err
}

func ReadAll() ([]*Schedule, error) {
	conn := newDb()
	defer conn.Close()

	var schedules []*Schedule
	err := newDb().
		Model(&Schedule{}).
		Order("start ASC").
		Select(&schedules)

	return schedules, err
}

func Update(s *Schedule) error {
	conn := newDb()
	defer conn.Close()

	_, err := conn.Model(s).Update(s)
	return err
}

func Delete(id int64) error {
	conn := newDb()
	defer conn.Close()

	_schedule, err := Read(id)
	if err != nil{
		return err
	}

	_, err = conn.
		Model(&_schedule).
		Delete(_schedule)

	return err
}
