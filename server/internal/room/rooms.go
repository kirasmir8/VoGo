package room

import (
	"fmt"
	"log/slog"
	"sync"
)

type Rooms struct {
	room map[string]*Room
	log  *slog.Logger
	mux  sync.Mutex
}

func NewRooms() *Rooms {
	return &Rooms{
		room: make(map[string]*Room),
		mux:  sync.Mutex{},
	}
}

func (r *Rooms) AddRoom(name string) error {
	const op = "internal.handlers"

	if _, ok := r.room[name]; ok {
		return fmt.Errorf("Комната %s уже существует: %w", name, op)
	}
	r.mux.Lock()
	r.room[name] = NewRoom()
	r.mux.Unlock()
	//TODO: придумать что-то с логгированием
	fmt.Println("Комната успешно создана", slog.String("name", name))
	return nil
}
