package wshandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/service"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait       = 10 * time.Second
	pongWait        = 60 * time.Second
	pingPeriod      = (pongWait * 9) / 10
	maxMessageSize  = 512
	websocketBuffer = 1024
)

type subscription struct {
	conn *connection
	room string
}

type connection struct {
	ws   *websocket.Conn
	send chan []byte
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (s *subscription) readPump() {
	c := s.conn
	defer func() {
		H.unregister <- *s
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { 
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.ws.ReadMessage()
		fmt.Println("test msg", msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("read error: %v", err)
			}
			break
		}
		fmt.Println("Received message from client:", string(msg))

		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		m := message{s.room, msg}
		H.broadcast <- m
	}
}

type stockWebsocket struct {
	stockService service.StockService
}

func NewStockWebsocket(stockService service.StockService) stockWebsocket {
	return stockWebsocket{stockService}
}

func (s *subscription) writePump(h stockWebsocket) {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	fmt.Println("send message")
	go func() {

	}()
	for {
		select {
		case _, ok := <-c.send:
			if !ok {
				return
			}

			collection, err := h.stockService.GetStockCollection(s.room)
			if err != nil {
				fmt.Printf("error %s", err)
			}
			m := map[string]interface{}{
				"time":       time.Now().Format("15:04:05 | 2006-01-02"),
				"room":       s.room,
				"collection": collection,
			}
			jsonData, err := json.Marshal(m)
			if err != nil {
				log.Fatal(err)
			}

			if err := c.write(websocket.TextMessage, jsonData); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (h stockWebsocket) ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("this error", err)
	}

	queryValues := r.URL.Query()
	roomId := queryValues.Get("roomId")
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{ws, make(chan []byte, 256)}
	s := subscription{c, roomId}
	s.room = strings.Trim(s.room, " ")
	hub.register <- s
	go s.writePump(h)
	go s.readPump()

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			message := fmt.Sprintf("stock -> %s", roomId)
			s.conn.send <- []byte(message)
		}
	}()
}
