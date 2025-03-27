package room

import (
	"fmt"
	"gitlab.com/kirasmir2/vogo/server/internal/participant"
	"log/slog"
	"os/user"
	"sync"
)

type Room struct {
	log          *slog.Logger
	mux          *sync.Mutex
	Participants map[string]*participant.Participant
}

func NewRoom() *Room {
	return &Room{
		Participants: make(map[string]*participant.Participant),
		mux:          &sync.Mutex{},
	}
}

// GetAllParticipants - возвращает список участников данного лобби
func (r *Room) GetAllParticipants() []string {
	participants := make([]string, 0, len(r.Participants))
	for part := range r.Participants {
		participants = append(participants, part)
	}
	return participants
}

func (r *Room) AddParticipant(nameUser string, participants *participant.Participant) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	_, ok := r.Participants[nameUser]
	if ok {
		return fmt.Errorf("пользователь %s уже существует", nameUser)
	}
	// добавляем участника в комнату
	r.Participants[nameUser] = participants
	fmt.Printf("пользователь %s успешно подключился")
	return nil
}

func (r *Room) RemoveParticipant(userInfo *user.User) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	_, ok := r.Participants[userInfo.Name]
	if ok {
		delete(r.Participants, userInfo.Name)
		return nil
	}
	return fmt.Errorf("пользователь %s не найден", userInfo.Name)
}

// BroadCastMessage - широковещательное сообщение всем участникам в данной комнате
func (r *Room) BroadCastMessage(message []byte) {
	for name, participant := range r.Participants {
		err := participant.SendMessage(message)
		// ошибка означает проблемы с пользователем, удаляем его из комнаты
		if err != nil {
			delete(r.Participants, name)
		}
	}
}
