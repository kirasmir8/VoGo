package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"gitlab.com/kirasmir2/vogo/server/internal/model"
	"gitlab.com/kirasmir2/vogo/server/internal/room"
	"log/slog"
	"net/http"
)

type Controller struct {
	rooms    *room.ActiveRooms
	log      *slog.Logger
	upgrader *websocket.Upgrader
}

func NewController(log *slog.Logger) *Controller {
	return &Controller{
		rooms: room.NewRooms(),
		log:   log,
		// TODO: Увеличить буферы до 4096 при деплои
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{Message: "Комната создана"})
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
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.log.Error("Ошибка апгрейда соединения", "error", err)
		return
	}
	defer conn.Close()
	c.log.Info("Клиент подключён", "addr", conn.RemoteAddr().String(), "user", userName, "room", roomName)

	if err := c.rooms.AddParticipant(userName, roomName, conn); err != nil {
		c.log.Error("Ошибка добавления участника", "error", err, "user", userName, "room", roomName)
		//TODO: добавить ответное сообщение об имени
		return
	}

	//test ехо сообщение
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			c.log.Error(err.Error())
			break
		}
		c.rooms.Rooms[roomName].BroadCastMessage(message)
	}
}

func (c *Controller) GetRoomsHandler(w http.ResponseWriter, r *http.Request) {
	activeRooms := c.rooms.GetRooms()

	res, err := json.Marshal(activeRooms)
	if err != nil {
		c.log.Error("Ошибка сериализации комнат", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
