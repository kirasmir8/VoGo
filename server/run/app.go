package run

import (
	"gitlab.com/kirasmir2/vogo/server/internal/handlers"
	"gitlab.com/kirasmir2/vogo/server/internal/infrastructure/router"
	"gitlab.com/kirasmir2/vogo/server/internal/infrastructure/server"
	"log/slog"
	"os"
)

type App struct {
	log *slog.Logger
	srv *server.AudioServer
}

func NewApp() *App {
	return &App{}
}

//TODO: Конструктор

func (a *App) Start() {
	a.srv.MustStart()
}

func (a *App) Init() *App {
	// инициализация логгера
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	a.log = logger
	// инициализация сервера
	srv := server.NewServer("8080", logger)
	a.srv = srv
	// инициализация контроллера
	controller := handlers.NewController(logger)
	// инициализация роутера
	router := router.NewRout(controller)
	// настройка сервера
	a.srv.Server.Handler = router
	return a
}
