package client

import "github.com/gorilla/websocket"

// Participant - структура участников голосового чата
type Participant struct {
	Conn *websocket.Conn
}

// InitParticipant - инициализирует нового участника
func InitParticipant(conn *websocket.Conn) *Participant {
	return &Participant{conn}
}

// TODO: Сделать проверку соединения по прицнипу Ping Pong
// CheckConnection - проверяет соединение с участником
func (p *Participant) CheckConnection() {
	p.Conn.PingHandler()
}
