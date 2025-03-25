package rout

import (
	"github.com/go-chi/chi"
	"gitlab.com/kirasmir2/vogo/server/internal/handlers"
)

func NewRout(controller *handlers.Controller) chi.Router {
	r := chi.NewRouter()
	r.Post("/room/{name}", controller.CreateRoomRequestHandler)
	r.Get("/test", controller.HandleWebSocket)

	return r
}
