package tailer

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/fsnotify/fsnotify"
)

func WatchLog(filename string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("Error initialzing watcher: %w \n", err)
	}
	defer watcher.Close()

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Error opening the file at %s : %w", filename, err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("Unable to get file stats: %w", err)
	}
	offset := info.Size()
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = watcher.Add(filename)
	if err != nil {
		return fmt.Errorf("Error watching the file %w", err)
	}

	reader := bufio.NewReader(file)

	fmt.Printf("Tailer active: Watching %s from offset %d \n", filename, offset)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			if event.Has(fsnotify.Write) {
				for {
					line, err := reader.ReadString('\n')

					if len(line) > 0 {
						fmt.Printf("New log line: %s", line)
					}

					if err != nil {
						if err == io.EOF {
							break
						}

						return fmt.Errorf("error reading the file %w", err)
					}
				}
			}

			if event.Has(fsnotify.Rename) || event.Has(fsnotify.Remove) {
				fmt.Println("Log rotation detected.")

				file.Close()

				var newFile *os.File
				for {
					newFile, err = os.Open(filename)
					if err == nil {
						break
					}
				}

				file = newFile
				reader = bufio.NewReader(file)

				watcher.Remove(event.Name)
				watcher.Add(filename)

			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			return fmt.Errorf("watcher error %w", err)
		}
	}
}
