package api

import (
	"encoding/json"
	"io"
	"net/http"

	"compton-scat/internal/runbook"
)

type Server struct {
	mux  *http.ServeMux
	addr string
	book *runbook.Book
}

func New(addr string) *Server {
	s := &Server{
		mux:  http.NewServeMux(),
		addr: addr,
		book: runbook.NewBook(64),
	}
	s.routes()
	return s
}

func Serve(addr string) error {
	return New(addr).ListenAndServe()
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) Book() *runbook.Book {
	return s.book
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.addr, s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/kinematics", s.handleKinematics)
	s.mux.HandleFunc("/api/section", s.handleSection)
	s.mux.HandleFunc("/api/trend", s.handleTrend)
	s.mux.HandleFunc("/api/history", s.handleHistory)
	s.mux.HandleFunc("/api/health", s.handleHealth)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET only")
		return
	}
	writeJSON(w, http.StatusOK, s.book.List())
}

func readJSON(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return errEmptyBody
	}
	return json.Unmarshal(body, v)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

var errEmptyBody = &json.SyntaxError{}
