package p2p

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Validator struct {
	conn *websocket.Conn
}

func NewValidator(addr string) (*Validator, error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}
	v := &Validator{conn: conn}
	go v.handleMessages()
	return v, nil
}

func (v *Validator) handleMessages() {
	defer v.conn.Close()
	for {
		messageType, p, err := v.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		if messageType == websocket.TextMessage {
			log.Println("Validator - Received:", string(p))
		}
	}
}

func (v *Validator) Ping() error {
	err := v.conn.WriteMessage(websocket.TextMessage, []byte("Ping"))
	if err != nil {
		return err
	}

	// Since we're handling messages in a separate goroutine,
	// we don't need to read the response here. The handleMessages method will log it.
	return nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Sequencer struct {
	server  *http.Server
	clients []*websocket.Conn
}

func NewSequencer(addr string) *Sequencer {
	s := &Sequencer{}
	s.server = &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println("Upgrade error:", err)
				return
			}
			s.clients = append(s.clients, conn)
			go s.handleClient(conn)
		}),
	}

	return s
}

func (s *Sequencer) handleClient(conn *websocket.Conn) {
	defer conn.Close()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		// If the message type is text and the content is "Ping"
		// Respond with a "Pong"
		if messageType == websocket.TextMessage && string(p) == "Ping" {
			conn.WriteMessage(websocket.TextMessage, []byte("Pong"))
			log.Println("Sequencer - Received:", string(p))
		}
	}
}

func (s *Sequencer) BroadcastHello() {
	for _, client := range s.clients {
		client.WriteMessage(websocket.TextMessage, []byte("Hello"))
	}
}

func (s *Sequencer) Start() {
	log.Fatal(s.server.ListenAndServe())
}
