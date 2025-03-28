package server

import (
	"go.uber.org/zap"
	"net/http"
)

// AudioServer - структура сервера
type AudioServer struct {
	Server *http.Server
	log    *zap.Logger
}

// NewServer - создание нового сервера
func NewServer(port string, logger *zap.Logger) *AudioServer {
	return &AudioServer{
		Server: &http.Server{
			Addr: "0.0.0.0:" + port,
		},
		log: logger,
	}
}

func (as *AudioServer) MustStart() {
	as.log.Info("Starting server", zap.String("address", as.Server.Addr))

	err := as.Server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
