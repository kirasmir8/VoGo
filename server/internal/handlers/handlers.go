package handlers

import (
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"gitlab.com/kirasmir2/vogo/server/internal/infrastructure/api"
	"gitlab.com/kirasmir2/vogo/server/internal/model"
	"gitlab.com/kirasmir2/vogo/server/internal/room"
	"go.uber.org/zap"
	"net/http"
)

type Controller struct {
	rooms    *room.ActiveRooms
	log      *zap.Logger
	upgrader *websocket.Upgrader
}

func NewController(log *zap.Logger) *Controller {
	return &Controller{
		rooms: room.NewRooms(log),
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
	//TODO: Обрезать пробелы справа и слева
	name := chi.URLParam(r, "name")
	if name == "" {
		api.StatusMessageResponse(
			w,
			http.StatusBadRequest,
			model.Response{Message: "название комнаты не указано"})
		return
	}

	err := c.rooms.AddRoom(name)
	if err != nil {
		api.StatusMessageResponse(
			w,
			http.StatusBadRequest,
			model.Response{Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	api.StatusMessageResponse(w, http.StatusOK, model.Response{Message: "Комната успешно создана"})
}

func (c *Controller) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// валидация запроса
	roomName := r.URL.Query().Get("room")
	userName := r.URL.Query().Get("name")
	if roomName == "" {
		api.StatusMessageResponse(
			w,
			http.StatusBadRequest,
			model.Response{Message: "Для подключения необходимо указать название комнаты"})
		return
	} else if userName == "" {
		api.StatusMessageResponse(
			w,
			http.StatusBadRequest,
			model.Response{Message: "Для подключения необходимо указать имя"})
		return
	}

	//Преобразование http-запроса в WebSocket
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.log.Error("Ошибка апгрейда соединения", zap.Error(err))
		return
	}
	defer conn.Close()
	c.log.Info("успешная установка соединения,",
		zap.String("room", roomName),
		zap.String("user", userName),
		zap.String("conn", conn.RemoteAddr().String()))

	if err := c.rooms.AddParticipant(userName, roomName, conn); err != nil {
		c.log.Error("Ошибка добавления участника",
			zap.Error(err),
			zap.String("room", roomName),
			zap.String("user", userName))
		return
	}

	//test ехо сообщение
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			c.log.Error(
				"ошибка отправки сообщения",
				zap.String("user", userName),
				zap.String("room", roomName),
				zap.Error(err))
			err = c.rooms.Rooms[roomName].RemoveParticipant(userName)
			if err != nil {
				c.log.Error(err.Error())
			}
			c.log.Info(
				"участник успешно удален",
				zap.String("user", userName),
				zap.String("room", roomName),
				zap.Error(err))
			break
		}
		c.rooms.Rooms[roomName].BroadCastMessage(message, userName)
	}
}

func (c *Controller) GetRoomsHandler(w http.ResponseWriter, r *http.Request) {
	activeRooms := c.rooms.GetRooms()
	api.StatusMessageResponse(w, http.StatusOK, activeRooms)
}
