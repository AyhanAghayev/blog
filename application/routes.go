package application

import (
	"anan/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	fs := http.FileServer(http.Dir("./static"))

	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Get("/", handler.ServeAbout)
	router.Get("/projects", handler.ServeProjects)
	router.Get("/post/{slug}", handler.ServePost)
	router.Get("/posts", handler.ServePosts)

	return router

}
