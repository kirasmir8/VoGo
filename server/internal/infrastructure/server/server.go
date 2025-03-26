package server

import (
	"log/slog"
	"net/http"
)

// AudioServer - структура сервера
type AudioServer struct {
	Server *http.Server
	log    *slog.Logger
}

// NewServer - создание нового сервера
func NewServer(port string, logger *slog.Logger) *AudioServer {
	return &AudioServer{
		Server: &http.Server{
			Addr: "0.0.0.0:" + port,
		},
		log: logger,
	}
}

func (as *AudioServer) MustStart() {
	as.log.Info("Starting server", slog.String("host", as.Server.Addr))
	err := as.Server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
