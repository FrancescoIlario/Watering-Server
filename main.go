package main

import (
	"fmt"
	"github.com/FrancescoIlario/Watering-Server/database"
	"github.com/FrancescoIlario/Watering-Server/rest-api/schedule-rest"
	"github.com/FrancescoIlario/Watering-Server/schedule"
	"github.com/FrancescoIlario/Watering-Server/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"time"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/schedule", schedule_rest.Routes())
		r.Get("/hello", func(resp http.ResponseWriter, req *http.Request) {
			_, err := resp.Write([]byte("Hello"))
			utils.PanicIf(err)
		})
		r.Get("/echo", Echo)
	})
	return router
}

func Echo(w http.ResponseWriter, r *http.Request) {
	var body []byte

	_, err := r.Body.Read(body)
	utils.PanicIf(err)

	_, err = w.Write(body)
	utils.PanicIf(err)

	w.WriteHeader(201)
}

func Hello(resp http.ResponseWriter, _ *http.Request) {
	_, err := resp.Write([]byte("Hello"))
	utils.PanicIf(err)
}

func startServer() {
	log.Println("configuring server")

	router := Routes()
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}
	err := chi.Walk(router, walkFunc)
	utils.PanicIf(err)

	log.Println("starting server")
	log.Fatal(http.ListenAndServe(":9999", router))
}

func testDb() {
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

		utils.PanicIf(database.Create(&newSchedule))
	}

	fmt.Printf("\n****\n\n")

	schedules, err := database.ReadAll()
	utils.PanicIf(err)

	for _, sched := range schedules {
		fmt.Println(sched.String())
	}
	fmt.Printf("\nTotally read: %v\n", len(schedules))
}

func main() {
	database.Initialize()
	startServer()
}
