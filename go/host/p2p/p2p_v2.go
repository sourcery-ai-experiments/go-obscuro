package p2p

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/obscuronet/go-obscuro/go/common/stopcontrol"
	"log"
	"net/http"
)

type Validator struct {
	conn          *websocket.Conn
	sequencerAddr string
	stopControl   *stopcontrol.StopControl
}

func NewValidator(addr string) *Validator {
	return &Validator{sequencerAddr: addr}
}

func (v *Validator) Start() error {
	conn, _, err := websocket.DefaultDialer.Dial(v.sequencerAddr, nil)
	if err != nil {
		return err
	}
	return v.connect(conn)

}

func (v *Validator) Stop() error {
	v.stopControl.Stop()

	return v.conn.Close()
}

func (v *Validator) SendMessage(msg []byte) error {
	return sendMessage(v.conn, msg)
}

func (v *Validator) connect(conn *websocket.Conn) error {
	// todo makes sure there's no handler in-flight
	v.conn = conn

	go v.messageHandler()

	return nil

}

func (v *Validator) messageHandler() {
	defer func() {
		fmt.Println("Error happened, closing conn")
		err := v.conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	for {
		// todo handle connection failures
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

type Sequencer struct {
	server  *http.Server
	clients []*websocket.Conn
}

func NewSequencer(addr string) *Sequencer {
	s := &Sequencer{}

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	s.server = &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println("Upgrade error:", err)
				return
			}
			s.clients = append(s.clients, conn)
			go s.handleClients(conn)
		}),
	}

	return s
}

func (s *Sequencer) Broadcast(data []byte) {
	for _, client := range s.clients {
		c := client
		go func() {
			err := c.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
}

func (s *Sequencer) Start() {
	log.Fatal(s.server.ListenAndServe())
}

func (s *Sequencer) Stop() error {
	return s.server.Close()
}

func (s *Sequencer) handleClients(conn *websocket.Conn) {
	defer func() {
		fmt.Println("Error happened, closing conn")
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		// If the message type is text and the content is "Ping"
		// Respond with a "Pong"
		if messageType == websocket.TextMessage {
			conn.WriteMessage(websocket.TextMessage, []byte("Pong"))
			log.Println("Sequencer - Received:", string(p))
		}
	}
}

func sendMessage(conn *websocket.Conn, bytes []byte) error {
	return conn.WriteMessage(websocket.TextMessage, bytes)
}

//func encode(msg interface{}) ([]byte, error) {
//	var buffer bytes.Buffer
//	enc := gob.NewEncoder(&buffer)
//	if err := enc.Encode(msg); err != nil {
//		return nil, err
//	}
//	return buffer.Bytes(), nil
//}
//
//func decode(data []byte) (*Message, error) {
//	var msg Message
//
//	if len(data) == 0 {
//		return nil, NoDataErr
//	}
//
//	buffer := bytes.NewBuffer(data)
//	dec := gob.NewDecoder(buffer)
//	if err := dec.Decode(&msg); err != nil {
//		return &msg, err
//	}
//	return &msg, nil
//}
