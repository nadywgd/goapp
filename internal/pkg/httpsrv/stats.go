package httpsrv

import "log"

type sessionStats struct {
	id   string
	sent int
}

func (w *sessionStats) print() {
	log.Printf("session %s has received %d messages\n", w.id, w.sent)
}

func (w *sessionStats) inc() {
	w.sent++
}

func (s *Server) incStats(id string) {
	// Find and increment.
	if ws, exists := s.sessionStats[id]; exists {
		ws.inc()
	} else {
		s.sessionStats[id] = &sessionStats{id: id, sent: 1}
	}
}
