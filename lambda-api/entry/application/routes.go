package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) loadRoutes() {
	ApiPrefix := "/api/v1/"

	router := chi.NewRouter()

	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	})

	router.Post(ApiPrefix+"execute", a.ExecuteRequest)

	router.Post(ApiPrefix+"test", a.TestRequest)

	a.router = router
}
