package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ParthK7/GoStash/internal/tailer"
	"github.com/ParthK7/GoStash/internal/wal"
	"github.com/gorilla/websocket"
)

type Server struct {
	logger   *wal.Wal
	hub      *Hub
	upgrader websocket.Upgrader
}

func (s *Server) handleIngest(w http.ResponseWriter, req *http.Request) {
	//from req check if it is post and if yes then use the logger to write
	// error helper function in http package -> http.Error(w, message, code)
	if req.Method != http.MethodPost {
		http.Error(w, "Only POST methods are allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(req.Body)
	defer req.Body.Close() //Best practice -> important to close it to ensure that the socket is closed, telling the server its communicationn is done.
	if err != nil || len(body) == 0 {
		http.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}

	bodyStr := string(body)

	err = s.logger.Write(bodyStr)
	if err != nil {
		http.Error(w, "Could not log request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data logged successfully"))
}

func (s *Server) handleWs(w http.ResponseWriter, req *http.Request) {
	connection, err := s.upgrader.Upgrade(w, req, nil)
	if err != nil {
		http.Error(w, "Could not upgrade HTTP connection", http.StatusInternalServerError)
		return
	}
	defer func() {
		s.hub.unregister <- connection
		connection.Close()
	}()

	s.hub.register <- connection
	for {
		_, _, err := connection.ReadMessage()
		if err != nil {
			break
		}
	}
}

func main() {
	_ = os.MkdirAll("cmd/server/storage", 0755)

	logger, err := wal.NewWal("cmd/server/storage/active.log")
	if err != nil {
		log.Fatalf("Failed to initialize WAL: %v", err)
	}

	logPipe := make(chan string)

	hub := NewHub(logPipe)

	// start the consumer first
	go hub.Run()

	go func() {
		err := tailer.WatchLog("cmd/server/storage/active.log", logPipe)
		if err != nil {
			log.Printf("trailer error %v", err)
		}
	}()

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	myServer := Server{logger: logger, hub: hub, upgrader: upgrader}

	http.HandleFunc("/ingest", myServer.handleIngest)
	http.HandleFunc("/ws", myServer.handleWs)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.ListenAndServe(":8080", nil)

}
