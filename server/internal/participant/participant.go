package participant

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Participant - структура участников голосового чата
type Participant struct {
	Conn *websocket.Conn
	log  *zap.Logger
}

// InitParticipant - инициализирует нового участника
func InitParticipant(conn *websocket.Conn, log *zap.Logger) *Participant {
	return &Participant{conn, log}
}

// SendMessage отправляет сообщение пользователю
func (p *Participant) SendMessage(message []byte) error {
	//TODO: при деплое исправить на бинарный формат
	if err := p.Conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
		// в случае ошибки при отправке, прекращаем соединение (если это возможно)
		p.Conn.Close()
		return err
	}
	return nil
}

func (p *Participant) Close() {
	err := p.Conn.Close()
	if err != nil {
		p.log.Error(
			"при попытке закрытие соединения произошла ошибка",
			zap.String("user", p.Conn.RemoteAddr().String()),
			zap.Error(err))
	}
}
