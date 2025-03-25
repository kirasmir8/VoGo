package room

import (
	"fmt"
	"gitlab.com/kirasmir2/vogo/server/internal/client"
	"log/slog"
	"os/user"
	"sync"
)

type Room struct {
	log          *slog.Logger
	mux          sync.Mutex
	Participants map[string]client.Participant
}

func NewRoom() *Room {
	return &Room{
		Participants: make(map[string]client.Participant),
		mux:          sync.Mutex{},
	}
}

func (r *Room) AddParticipants(userInfo *user.User, participants client.Participant) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	_, ok := r.Participants[userInfo.Name]
	if ok {
		return fmt.Errorf("Пользователь %s уже существует", userInfo.Name)
	}
	// добавляем участника в комнату
	r.Participants[userInfo.Name] = participants
	return nil
}

func (r *Room) RemoveParticipants(userInfo *user.User) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	_, ok := r.Participants[userInfo.Name]
	if ok {
		delete(r.Participants, userInfo.Name)
		return nil
	}
	return fmt.Errorf("Пользователь %s не найден", userInfo.Name)
}
