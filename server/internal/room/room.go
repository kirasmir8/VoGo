package room

import (
	"fmt"
	"gitlab.com/kirasmir2/vogo/server/internal/participant"
	"go.uber.org/zap"
	"sync"
)

type Room struct {
	log          *zap.Logger
	mux          *sync.Mutex
	Participants map[string]*participant.Participant
}

func NewRoom(log *zap.Logger) *Room {
	return &Room{
		Participants: make(map[string]*participant.Participant),
		mux:          &sync.Mutex{},
		log:          log,
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
	return nil
}

func (r *Room) RemoveParticipant(name string) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	_, ok := r.Participants[name]
	if ok {
		delete(r.Participants, name)
		return nil
	}
	return fmt.Errorf("пользователь %s не найден", name)
}

// BroadCastMessage - широковещательное сообщение всем участникам в данной комнате
func (r *Room) BroadCastMessage(message []byte, participantName string) {
	for name, part := range r.Participants {
		// TODO: Потом исправить на проверку поулучше
		if name == participantName {
			continue
		}
		err := part.SendMessage(message)
		// ошибка означает проблемы с пользователем, удаляем его из комнаты
		if err != nil {
			delete(r.Participants, name)
		}
	}
}
