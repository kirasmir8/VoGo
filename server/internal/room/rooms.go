package room

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gitlab.com/kirasmir2/vogo/server/internal/participant"
	"log/slog"
	"sync"
)

type Rooms struct {
	Room map[string]*Room
	log  *slog.Logger
	mux  sync.Mutex
}

func NewRooms() *Rooms {
	return &Rooms{
		Room: make(map[string]*Room),
		mux:  sync.Mutex{},
	}
}

// AddRoom - добавляет комнату в общий список комнат
func (r *Rooms) AddRoom(name string) error {
	const op = "internal.handlers"

	if _, ok := r.Room[name]; ok {
		return fmt.Errorf("Комната %s уже существует: %w", name, op)
	}
	r.mux.Lock()
	r.Room[name] = NewRoom()
	r.mux.Unlock()
	//TODO: придумать что-то с логгированием
	fmt.Println("Комната успешно создана", slog.String("name", name))
	return nil
}

// AddParticipant - добавляет участника в комнату
func (r *Rooms) AddParticipant(participantName string, roomName string, conn *websocket.Conn) error {
	if _, ok := r.Room[roomName]; !ok {
		return fmt.Errorf("комнаты %s не существует", roomName)
	} else if _, ok := r.Room[roomName].Participants[participantName]; ok {
		return fmt.Errorf("участник %s уже существует", roomName)
	}
	r.mux.Lock()
	r.Room[roomName].Participants[participantName] = participant.InitParticipant(conn)
	r.mux.Unlock()
	fmt.Printf("Участник %s успешно зашёл в комнату", participantName)
	return nil
}
