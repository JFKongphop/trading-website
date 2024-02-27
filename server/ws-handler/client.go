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

type stockWebsocket struct {
	stockService service.StockService
}

func NewStockWebsocket(stockService service.StockService) stockWebsocket {
	return stockWebsocket{stockService}
}

type connection struct {
	ws   *websocket.Conn
	send chan []byte
}

type subscription struct {
	conn *connection
	room string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

const (
	writeWait       = 10 * time.Second
	pongWait        = 60 * time.Second
	pingPeriod      = (pongWait * 9) / 10
	maxMessageSize  = 512
	websocketBuffer = 1024
)

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

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (s *subscription) writePrice(h stockWebsocket) {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case _, ok := <-c.send:
			if !ok {
				return
			}

			room := strings.Split(s.room, "-")[1]
			price, err := h.stockService.GetStockPrice(room)
			if err != nil {
				log.Printf("error %s", err)
			}

			m := map[string]interface{}{
				"time":  time.Now().Format("15:04:05 | 2006-01-02"),
				"room":  s.room,
				"pirce": price,
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

func (s *subscription) writeTransaction(h stockWebsocket) {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case _, ok := <-c.send:
			if !ok {
				return
			}

			room := strings.Split(s.room, "-")[1]
			tx, err := h.stockService.GetStockHistory(room)
			if err != nil {
				log.Printf("error %s", err)
			}

			m := map[string]interface{}{
				"time": time.Now().Format("15:04:05 | 2006-01-02"),
				"room": room,
				"tx":   tx,
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

func (s *subscription) writeGraph(h stockWebsocket) {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case _, ok := <-c.send:
			if !ok {
				return
			}

			room := strings.Split(s.room, "-")[1]
			graph, err := h.stockService.GetStockGraph(room)
			if err != nil {
				log.Printf("error %s", err)
			}

			m := map[string]interface{}{
				"graph": graph,
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

func (h stockWebsocket) ServePriceWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("this error", err)
	}

	queryValues := r.URL.Query()
	roomId := queryValues.Get("stockId")
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{ws, make(chan []byte, 256)}
	s := subscription{c, roomId}

	s.room = fmt.Sprintf("price-%s", strings.Trim(s.room, " "))
	hub.register <- s
	go s.writePrice(h)
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

func (h stockWebsocket) ServeTransactionWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("this error", err)
	}

	queryValues := r.URL.Query()
	roomId := queryValues.Get("stockId")
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{ws, make(chan []byte, 256)}
	s := subscription{c, roomId}

	s.room = fmt.Sprintf("tx-%s", strings.Trim(s.room, " "))
	hub.register <- s
	go s.writeTransaction(h)
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

func (h stockWebsocket) ServeGraphWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("this error", err)
	}

	queryValues := r.URL.Query()
	roomId := queryValues.Get("stockId")
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{ws, make(chan []byte, 256)}
	s := subscription{c, roomId}

	s.room = fmt.Sprintf("tx-%s", strings.Trim(s.room, " "))
	hub.register <- s
	go s.writeGraph(h)
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
