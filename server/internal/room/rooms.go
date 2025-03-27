package room

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gitlab.com/kirasmir2/vogo/server/internal/participant"
	"log/slog"
	"sync"
)

type ActiveRooms struct {
	Rooms map[string]*Room
	log   *slog.Logger
	mux   sync.Mutex
}

func NewRooms() *ActiveRooms {
	return &ActiveRooms{
		Rooms: make(map[string]*Room),
		mux:   sync.Mutex{},
	}
}

// AddRoom - добавляет комнату в общий список комнат
func (r *ActiveRooms) AddRoom(name string) error {
	const op = "internal.handlers"

	if _, ok := r.Rooms[name]; ok {
		return fmt.Errorf("комната %s уже существует: %s", name, op)
	}
	r.mux.Lock()
	r.Rooms[name] = NewRoom()
	r.mux.Unlock()
	//TODO: придумать что-то с логгированием
	fmt.Println("Комната успешно создана", slog.String("name", name))
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
	p := participant.InitParticipant(conn)
	if err := room.AddParticipant(participantName, p); err != nil {
		return err
	}
	fmt.Println("Участник добавлен", "user", participantName, "room", roomName)
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
