package schedule_rest

import (
	"github.com/filario/watering-server/database"
	"github.com/filario/watering-server/schedule"
	"github.com/filario/watering-server/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

func Routes() *chi.Mux {
	router := chi.NewMux()

	router.Get("/", ReadAll)
	router.Post("/", Create)
	router.Get("/{scheduleID}", Read)
	router.Put("/{scheduleID}", Update)
	router.Delete("/{scheduleID}", Delete)

	return router
}

func Create(w http.ResponseWriter, r *http.Request) {
	_schedule := schedule.Schedule{}
	err := render.DecodeJSON(r.Body, _schedule)
	utils.PanicIf(err) // TODO: handle this

	err = database.Create(&_schedule)
	utils.PanicIf(err)

	w.WriteHeader(201)
}

func ReadAll(w http.ResponseWriter, rq *http.Request) {
	schedules, err := database.ReadAll()
	utils.PanicIf(err)

	render.JSON(w, rq, schedules)
}

// /schedule/{scheduleId}
func Read(w http.ResponseWriter, r *http.Request) {
	scheduleID := chi.URLParam(r, "scheduleID")

	scheduleIDInt, err := strconv.ParseInt(scheduleID, 10, 64)
	utils.PanicIf(err) // TODO: handle this case and return error to client

	_schedule, err := database.Read(scheduleIDInt)
	utils.PanicIf(err)
	render.JSON(w, r, _schedule)
}

// /schedule/{scheduleId}
func Update(w http.ResponseWriter, r *http.Request) {
	scheduleID := chi.URLParam(r, "scheduleID")
	scheduleIDInt, err := strconv.ParseInt(scheduleID, 10, 64)
	utils.PanicIf(err) // TODO: handle this case and return error to client

	_schedule := schedule.Schedule{}
	err = render.DecodeJSON(r.Body, _schedule)
	utils.PanicIf(err) // TODO: handle this
	_schedule.Id = scheduleIDInt

	err = database.Update(&_schedule)
	utils.PanicIf(err)

	w.WriteHeader(200)
}

// /schedule/{scheduleId}
func Delete(w http.ResponseWriter, r *http.Request) {
	scheduleID := chi.URLParam(r, "scheduleID")

	scheduleIDInt, err := strconv.ParseInt(scheduleID, 10, 64)
	utils.PanicIf(err) // TODO: handle this case and return error to client

	err = database.Delete(scheduleIDInt)
	utils.PanicIf(err)

	w.WriteHeader(200)
}