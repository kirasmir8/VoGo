package participant

import (
	"fmt"
	"github.com/gorilla/websocket"
)

// Participant - структура участников голосового чата
type Participant struct {
	Conn *websocket.Conn
}

// InitParticipant - инициализирует нового участника
func InitParticipant(conn *websocket.Conn) *Participant {
	return &Participant{conn}
}

// SendMessage отправляет сообщение пользователю
func (p *Participant) SendMessage(message []byte) error {
	//TODO: при деплое исправить на бинарный формат
	if err := p.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
		// в случае ошибки при отправке, прекращаем соединение (если это возможно)
		p.Close()
		return err
	}
	return nil
}

func (p *Participant) Close() {
	err := p.Conn.Close()
	if err != nil {
		fmt.Printf("при попытке закрытие соединения произошла ошибка, участник,\nerror: %v\nclient: %s", err, p.Conn.RemoteAddr().String())
	}
}
