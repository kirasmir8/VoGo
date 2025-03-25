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
func NewServer(port string) *AudioServer {
	return &AudioServer{
		Server: &http.Server{
			Addr: ":" + port,
		},
	}
}

func (as *AudioServer) MustStart() {
	as.log.Info("Starting server", slog.String("host", as.server.Addr))
	err := as.Server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
