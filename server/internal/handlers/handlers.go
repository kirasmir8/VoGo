package handlers

import (
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"gitlab.com/kirasmir2/vogo/server/internal/room"
	"log/slog"
	"net/http"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаются все источники
	},
}

type Controller struct {
	rooms *room.Rooms
	log   *slog.Logger
}

func NewController(log *slog.Logger) *Controller {
	return &Controller{rooms: room.NewRooms(), log: log}
}

func (c *Controller) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
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

func (c *Controller) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// валидация запроса
	roomName := r.URL.Query().Get("room")
	userName := r.URL.Query().Get("name")
	if roomName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Для подключения необходимо указать название комнаты"))
		return
	} else if userName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Для подключения заполните имя пользователя"))
		return
	}

	//Преобразование http-запроса в WebSocket
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		c.log.Error(err.Error())
		return
	}
	defer conn.Close()
	c.log.Info("Подключение клиента", conn.RemoteAddr().String())

	err = c.rooms.AddParticipant(userName, roomName, conn)
	if err != nil {
		c.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	//test ехо сообщение
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			c.log.Error(err.Error())
			break
		}
		c.rooms.Room[roomName].BroadCastMessage(message)
	}
}
