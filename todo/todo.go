package todo

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

type Todo struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func Routes() *chi.Mux {
	router := chi.NewMux()

	router.Get("/", GetAllTodos)
	return router
}

func GetAllTodos(w http.ResponseWriter, rq *http.Request) {
	todos := []Todo{
		{
			"first",
			"content",
		},
	}
	render.JSON(w, rq, todos)


}
