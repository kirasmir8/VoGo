package room

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gitlab.com/kirasmir2/vogo/server/internal/participant"
	"go.uber.org/zap"
	"sync"
)

type ActiveRooms struct {
	Rooms map[string]*Room
	log   *zap.Logger
	mux   sync.Mutex
}

func NewRooms(log *zap.Logger) *ActiveRooms {
	return &ActiveRooms{
		Rooms: make(map[string]*Room),
		mux:   sync.Mutex{},
		log:   log,
	}
}

// AddRoom - добавляет комнату в общий список комнат
func (r *ActiveRooms) AddRoom(name string) error {
	if _, ok := r.Rooms[name]; ok {
		return fmt.Errorf("комната %s уже существует", name)
	}
	r.mux.Lock()
	r.Rooms[name] = NewRoom(r.log)
	r.mux.Unlock()

	r.log.Info("комната успешно создана", zap.String("name", name))
	return nil
}

// AddParticipant - добавляет участника в комнату
func (r *ActiveRooms) AddParticipant(participantName string, roomName string, conn *websocket.Conn) error {
	r.mux.Lock()
	room, ok := r.Rooms[roomName]
	r.mux.Unlock()
	if !ok {
		return fmt.Errorf("комнаты %s не существует", roomName)
	}
	p := participant.InitParticipant(conn, r.log)
	if err := room.AddParticipant(participantName, p); err != nil {
		return err
	}
	r.log.Info("участник добавлен", zap.String("room", roomName), zap.String("participant", participantName))
	return nil
}

// GetRooms - возвращает список активных комнат
func (r *ActiveRooms) GetRooms() []string {
	activeRooms := make([]string, 0, len(r.Rooms))
	for room, _ := range r.Rooms {
		activeRooms = append(activeRooms, room)
	}
	return activeRooms
}
