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
	if r.Method != http.MethodPost { // если метод не post, то выдаст ошибку
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	req := struct { // в структуру локальную записываем эти данные, присланные в post
		CandidateID uint64 `json:"candidate_id"`
		Passport    string `json:"passport"`
	}{}
	// get request body вычитываем тело post
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if err := json.Unmarshal(raw, &req); err != nil { // раскоживание данных, чтобы замаршалить в
		// структуру выше. в роу передаём данне и указываем структуру, куда будем маршалить
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	// validate
	if len(req.Passport) == 0 || req.CandidateID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.mu.Lock() // делаем счётчик, чтобы никто первее нас не положил данные в структуру
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

	vals := r.URL.Query() // куери параметры для запросы статы на конкретного кандидата
	if vals.Has(candID) { // если в валс есть кандидат айди, то он нам его возвращает
		id, err := strconv.ParseUint(vals.Get(candID), 10, 64) // всё что приходит в куери
		// параметрах представленно в виде строки. Нам надо законвертить строку в юинт
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		s.mu.RLock()             // лочим мьютекс, чтобы считать мапу для её отправки
		candStats := s.stats[id] // айди законвертированная куери
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

func New() srv { // конструктор указывающийся в майне
	return srv{
		mu:    &sync.RWMutex{},
		stats: make(map[uint64]uint),
	}
}
