package handlers

import (
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"gitlab.com/kirasmir2/vogo/server/internal/client"
	"gitlab.com/kirasmir2/vogo/server/internal/room"
	"log/slog"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаются все источники
	},
}

type Controller struct {
	rooms       *room.Rooms
	log         *slog.Logger
	participant *client.Participant
}

func NewController(participant *client.Participant, rooms *room.Rooms) *Controller {
	return &Controller{participant: participant, rooms: rooms}
}

func (c *Controller) CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "{name}")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("укажите имя комнаты"))
		return
	}

	err := c.rooms.AddRoom(name)
	if err != nil {
		c.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}
