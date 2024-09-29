package service

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"
)

type srv struct {
	mu    *sync.RWMutex
	stats map[uint64]uint
}

func (s *srv) Vote(w http.ResponseWriter, r *http.Request) {

	// check request method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	w.Write([]byte("Hello Server1657"))

	req := struct {
		CandidateID uint64 `json:"candidate_id"`
		Passport    string `json:"passport"`
	}{}
	// get request body
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if err := json.Unmarshal(raw, &req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	// validate
	if len(req.Passport) == 0 || req.CandidateID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.stats[req.CandidateID]++
	s.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func (s *srv) Stats(w http.ResponseWriter, r *http.Request) {
	const candID = "candidate_id"

	// check request method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	vals := r.URL.Query()
	if vals.Has(candID) {
		id, err := strconv.ParseUint(vals.Get(candID), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		s.mu.RLock()
		candStats := s.stats[id]
		s.mu.RUnlock()

		raw, err := json.Marshal(candStats)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.Write(raw)
		return
	}

	s.mu.RLock()
	stats := s.stats
	s.mu.RUnlock()

	raw, err := json.Marshal(stats)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Write(raw)
}

func New() srv {
	return srv{
		mu:    &sync.RWMutex{},
		stats: make(map[uint64]uint),
	}
}
