package database

import (
	"github.com/FrancescoIlario/Watering-Server/schedule"
	"github.com/FrancescoIlario/Watering-Server/utils"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"log"
)

type ICqlDbSchedule interface {
	connect() *gocql.Session
}

type CqlDbSchedule struct {}

func (_ *CqlDbSchedule) connect() *gocql.Session {
	// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Port = 9042
	cluster.Keyspace = "watering"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()

	utils.PanicIf(err)
	return session
}

func (db *CqlDbSchedule) ReadAll() []*schedule.Schedule {
	session := db.connect()
	defer session.Close()

	stmt, _ := qb.Select("schedule").Where(qb.Eq("deleted")).OrderBy("start", qb.DESC).ToCql()
	log.Println(stmt)

	var schedules []*schedule.Schedule
	if err := gocqlx.Select(&schedules, session.Query(stmt)); err != nil {
		utils.PanicIf(err)
	}
	return schedules
}

func (db *CqlDbSchedule) Read(id int) []*schedule.Schedule {
	session := db.connect()
	defer session.Close()

	stmt, names := qb.Select("schedule").Where(qb.Eq("id")).ToCql()
	log.Println(stmt)

	var schedules []*schedule.Schedule
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M { "id": id })
	if err := q.SelectRelease(&schedules); err != nil {
		log.Fatal(err)
	}
	return schedules
}

func (db *CqlDbSchedule) Create(schedule schedule.Schedule) {
	session := db.connect()
	defer session.Close()

	stmt, names := qb.Insert("schedule").Columns("id", "start", "end").ToCql()
	log.Println(stmt)

	//sid, err := uuid.NewV1()
 	//utils.PanicIf(err)
	//schedule.Id = sid

	q := gocqlx.Query(session.Query(stmt), names).BindStruct(schedule)
	if err := q.ExecRelease(); err != nil {
		log.Fatal(err.Error())
	}
}