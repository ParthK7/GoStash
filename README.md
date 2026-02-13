# GoStash

GoStash is a high-performance, real-time log ingestion and distribution engine built in Go. It is designed to bridge the gap between high-velocity data producers and real-time observability dashboards by isolating disk I/O from network delivery.

The system guarantees that every incoming log is safely persisted to a sequential, append-only file before acknowledgment, while simultaneously broadcasting updates to connected clients via WebSockets.

## Core Architecture

- **Sequential Persistence Engine**: Utilizes a mutex-protected file writer to handle concurrent ingestion. Logs are committed to a `active.log` file using synchronous write patterns to ensure data integrity.
- **Automated Lifecycle Management**: Features built-in log rotation logic that archives active logs based on configurable size thresholds, preventing disk exhaustion.
- **Asynchronous Fan-out**: Implements a Hub-and-Spoke distribution model. A background tailing process monitors file-system events via `fsnotify`, pushing new data through Go channels to a WebSocket broadcaster.
- **Resource Isolation**: The ingestion path (Disk I/O) is strictly decoupled from the delivery path (Network), ensuring that slow frontend consumers do not create backpressure on log producers.

## Performance Metrics

Validated through high-concurrency load testing:
* **Throughput**: 4,100+ Requests Per Second (RPS)
* **Average Ingestion Latency**: 241µs
* **Concurrency**: Sustained performance under 50+ simultaneous worker threads.

## Technology Stack

- **Language**: Go (Golang)
- **Concurrency**: Goroutines, Channels, and Mutexes
- **Networking**: Gorilla WebSockets
- **Event Tracking**: fsnotify (OS-level file watching)
- **Persistence**: Sequential File I/O with manual rotation

## Project Structure

- `main.go`: Entry point and HTTP handler coordination.
- `hub.go`: Manages WebSocket client state and data broadcasting.
- `tailer.go`: Monitors the persistent log file for real-time updates.
- `wal.go`: Handles thread-safe appends to the sequential log file.
- `index.html`: Real-time visualization dashboard.

## Getting Started

### Prerequisites
- Go 1.20 or higher

### Installation & Execution
1. Clone the repository:
   ```bash
   git clone [https://github.com/your-username/gostash.git](https://github.com/your-username/gostash.git)
   cd gostash
   ```

2. Start the server:
   ```bash
   go run cmd/server/*.go
   ```

3. Open the dashboard
- Navigate to http://localhost:8080 in your browser.

Ingesting Logs
Send a log entry via a POST request using curl:

```bash
curl -X POST -d "System status: nominal" http://localhost:8080/ingest
```
