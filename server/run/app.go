package run

import (
	"gitlab.com/kirasmir2/vogo/server/internal/server"
	"log/slog"
	"net/http"
	"os"
)

type App struct {
	log *slog.Logger
	srv *http.Server
}

//TODO: Конструктор

func (a *App) Init() {
	// инициализация логгера
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	a.log = logger
	// инициализация сервера
	srv := server.NewServer("8080")
	a.srv = srv.Server
}
