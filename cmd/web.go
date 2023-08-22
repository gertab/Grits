package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"phi/parser"
	"phi/process"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8081", "http service address")

// func serveHome(w http.ResponseWriter, r *http.Request) {
// 	log.Println(r.URL)
// 	if r.URL.Path != "/" {
// 		http.Error(w, "Not found", http.StatusNotFound)
// 		return
// 	}
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	http.ServeFile(w, r, "./cmd/home.html")
// }

// setupAPI will start all Routes and their Handlers
func setupAPI() {
	hub := newHub()
	go hub.run()
	// http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			// case message := <-h.broadcast:
			// 	for client := range h.clients {
			// 		select {
			// 		case client.send <- message:
			// 		default:
			// 			close(client.send)
			// 			delete(h.clients, client)
			// 		}
			// 	}
		}
	}
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// todo remove for security (CORS)
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// re            *process.RuntimeEnvironment
	subscriberInfo *process.SubscriberInfo
	// processes := parser.ParseString(program)
	// process.InitializeProcesses(processes)
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.handleRequest(string(message))
		// c.hub.broadcast <- message
	}
}

func (c *Client) handleMonitorUpdates() {
	for {
		if c.subscriberInfo == nil {
			log.Println("processUpdate is nil")
			return
		}

		log.Println("Listening for process updates")

		updatedProcesses := <-c.subscriberInfo.ProcessesStringChan

		reply := &ReplyMessage{Type: "processes_updated", Payload: updatedProcesses}
		log.Println("RECEIVED process updates")

		c.send <- []byte("NEWWWWWW")
		// e.SetEscapeHTML(false)
		// e.Encode(statement)
		// reply_json, err := json.Marshal(reply)
		reply_json, err := reply.JSON()

		if err == nil {
			c.send <- []byte(reply_json)
		}
		// for _, p := range updatedProcesses {
		// 	c.send <- []byte(json.Marshal(updatedProcesses))
		// }
	}
}

type ReplyMessage struct {
	Type    string                `json:"type"`
	Payload []process.ProcessInfo `json:"payload"`
}

func (t *ReplyMessage) JSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

type RequestMessage struct {
	Type             string `json:"type"`
	ProgramToCompile string `json:"program_to_compile"`
}

func (c *Client) handleRequest(message string) {
	request := RequestMessage{}

	json.Unmarshal([]byte(message), &request)

	log.Println(request.Type)
	log.Println(request.ProgramToCompile)

	if request.Type == "compile_program" {
		log.Println("compiling program")

		c.subscriberInfo = process.NewSubscriberInfo()

		go c.handleMonitorUpdates()

		processes := parser.ParseString(request.ProgramToCompile)
		process.InitializeProcesses(processes, c.subscriberInfo)

		// c.send <- []byte(processes)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// // Add queued chat messages to the current websocket message.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write(newline)
			// 	w.Write(<-c.send)
			// }

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("err")
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
