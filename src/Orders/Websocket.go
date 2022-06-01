package Orders

import (
	"log"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type request struct {
	Op      string `json:"op"`
	Channel string `json:"channel"`
	Market  string `json:"market"`
}

func WebSocket(c chan []byte, channels string, symbols string) {

	conn, _, err := websocket.DefaultDialer.Dial("wss://ftx.com/ws/", nil)

	if err != nil {
		log.Println(err)
		return
	}

	err = subscribe(conn, channels, symbols)

	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	for {

		_, msg, err := conn.ReadMessage()

		if err != nil {
			log.Println("Error Reading Message: ", err)
		}

		// fmt.Println(string(msg))

		c <- msg

	}

}

func subscribe(conn *websocket.Conn, channels string, symbols string) error {

	err := conn.WriteJSON(&request{
		Op:      "subscribe",
		Channel: channels,
		Market:  symbols,
	})

	if err != nil {
		return err
	}

	err = conn.WriteJSON(&request{
		Op:      "subscribe",
		Channel: channels,
	})

	if err != nil {
		return err
	}

	return nil

}
