package httpsrv

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"goapp/internal/pkg/watcher"
	"goapp/pkg/util"

	"github.com/gorilla/websocket"
)

func (s *Server) handlerHexWebSocket(w http.ResponseWriter, r *http.Request) {
	// Create and start a watcher.
	var watch = watcher.New()
	if err := watch.Start(); err != nil {
		s.error(w, http.StatusInternalServerError, fmt.Errorf("failed to start watcher: %w", err))
		return
	}
	defer watch.Stop()

	s.addWatcher(watch)
	defer s.removeWatcher(watch)

	// Start WebSocket connection.
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origin == "http://localhost:8080"
		},
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.error(w, http.StatusInternalServerError, fmt.Errorf("failed to upgrade connection: %w", err))
		return
	}
	defer func() { _ = c.Close() }()

	log.Printf("websocket started for hex watcher %s\n", watch.GetWatcherId())
	defer func() {
		log.Printf("websocket stopped for hex watcher %s\n", watch.GetWatcherId())
	}()

	// Read done.
	readDoneCh := make(chan struct{})

	// All done.
	doneCh := make(chan struct{})
	defer close(doneCh)

	go func() {
		defer close(readDoneCh)
		for {
			select {
			default:
				_, p, err := c.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
						log.Printf("failed to read message: %v\n", err)
					}
					return
				}

				var m watcher.CounterReset
				if err := json.Unmarshal(p, &m); err != nil {
					log.Printf("failed to unmarshal message: %v\n", err)
					continue
				}

				// Validate the action
				if m.Action != "reset" {
					log.Printf("invalid action: %s\n", m.Action)
					continue
				}

				log.Printf("Resetting counter with value: %d\n", m.Value)
				watch.ResetCounter()

			case <-doneCh:
				return
			case <-s.quitChannel:
				return
			}
		}
	}()

	for {
		select {
		case cv := <-watch.Recv():
			// Marshal cv to JSON
			data, err := json.Marshal(cv)
			if err != nil {
				log.Printf("failed to marshal counter value: %v\n", err)
				return
			}

			// Convert the JSON to a map
			var cvMap map[string]interface{}
			if err := json.Unmarshal(data, &cvMap); err != nil {
				log.Printf("failed to unmarshal counter value: %v\n", err)
				return
			}

			// Generate the hex value and add to the map
			hexValue := util.RandHexString(10)
			cvMap["value"] = hexValue

			// Marshal the modified map back to JSON
			finalData, err := json.Marshal(cvMap)
			if err != nil {
				log.Printf("failed to marshal updated data: %v\n", err)
				return
			}

			// Send the message
			err = c.WriteMessage(websocket.TextMessage, finalData)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("failed to write message: %v\n", err)
				}
				return
			}
		case <-readDoneCh:
			return
		case <-s.quitChannel:
			return
		}
	}
}
