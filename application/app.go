package application

import (
	"fmt"
	"net/http"
)

type App struct {
	router http.Handler
}

func New() *App {
	return &App{
		router: loadRoutes(),
	}
}

func (a *App) Start() error {
	server := &http.Server{
		Handler: a.router,
		Addr:    ":3000",
	}

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error occured: %w", err)
	}

	return nil
}
