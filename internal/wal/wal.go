package wal

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Wal struct {
	filepath string
	logfile *os.File
	mu sync.Mutex
	lastRotatedTime time.Time
}

func NewWal(filepath string) (*Wal, error) {
	file, err := os.OpenFile(filepath, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return &Wal{filepath : filepath, logfile : file, lastRotatedTime : time.Now()}, nil
}


func (w *Wal)Write (content string) (error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	info, err := w.logfile.Stat()
	if err != nil {
		return err
	}
	
	// Change this to 10 * 1024 * 1024 later 
	if info.Size() > 5000 || time.Since(w.lastRotatedTime) > 2 * time.Minute {
		// rotation logic goes here 
		err = w.rotate()
		if err != nil {
			return err
		}
	}

	_, err = w.logfile.WriteString(content + "\n")
	if err != nil {
		return err
	}

	return nil
}

func (w *Wal) rotate() error {
	w.logfile.Close()

	newName := fmt.Sprintf("%s.%d", w.filepath, time.Now().UnixNano())
	err := os.Rename(w.filepath, newName)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(w.filepath, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	w.logfile = file
	w.lastRotatedTime = time.Now()
	return nil
}



