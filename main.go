package main

import (
	"fmt"
	"github.com/FrancescoIlario/Watering-Server/database"
	"github.com/FrancescoIlario/Watering-Server/rest-api/schedule-rest"
	simple_consumer "github.com/FrancescoIlario/Watering-Server/simple-consumer"
	simple_producer "github.com/FrancescoIlario/Watering-Server/simple-producer"
	"github.com/FrancescoIlario/Watering-Server/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
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

	bodyString := fmt.Sprintf("%v", body)
	err = simple_producer.Publish(&bodyString)
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

func main() {
	database.Initialize()

	go simple_consumer.StartConsumer(0)
	startServer()
}
