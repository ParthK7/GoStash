package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ParthK7/GoStash/internal/tailer"
	"github.com/ParthK7/GoStash/internal/wal"
)

type Server struct {
	logger *wal.Wal
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

func main() {
	_ = os.MkdirAll("storage", 0755)

	logger, err := wal.NewWal("storage/active.log")
	if err != nil {
		log.Fatalf("Failed to initialize WAL: %v", err)
	}

	go func() {
		err := tailer.WatchLog("storage/active.log")
		if err != nil {
			log.Printf("trailer error %v", err)
		}
	}()

	myServer := Server{logger: logger}

	http.HandleFunc("/ingest", myServer.handleIngest)

	http.ListenAndServe(":8080", nil)

}
