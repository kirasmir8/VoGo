package router

import (
	"github.com/go-chi/chi"
	"gitlab.com/kirasmir2/vogo/server/internal/handlers"
)

func NewRout(controller *handlers.Controller) chi.Router {
	r := chi.NewRouter()
	r.Post("/room/{name}", controller.CreateRoomHandler)
	r.Get("/room/connect", controller.HandleWebSocket)
	r.Get("/rooms", controller.GetRoomsHandler)

	return r
}
