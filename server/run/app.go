package run

import (
	"gitlab.com/kirasmir2/vogo/server/internal/handlers"
	zapLogger "gitlab.com/kirasmir2/vogo/server/internal/infrastructure/logger"
	"gitlab.com/kirasmir2/vogo/server/internal/infrastructure/router"
	"gitlab.com/kirasmir2/vogo/server/internal/infrastructure/server"
	"go.uber.org/zap"
)

type App struct {
	logger *zap.Logger
	srv    *server.AudioServer
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
	logger := zapLogger.MustNewLogger()
	a.logger = logger
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
