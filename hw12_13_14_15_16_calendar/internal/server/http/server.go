package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"io"
	"log"
	"net/http"
	"time"
)

type Logger interface {
	Sync()
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}

type Application interface {
}

type Server struct {
	server *http.Server
}

func NewServer(logger Logger, app Application, baseURL string, port int, readTimeout, writeTimeout time.Duration) *Server {
	r := chi.NewRouter()
	r.Use(LoggingMiddleware(logger))
	r.Route("/", func(r chi.Router) {
		r.Get("/", index)
	})

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", baseURL, port),
		Handler:      r,
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
	}
	return &Server{
		server: server,
	}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

func index(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}
	buf := make([]byte, r.ContentLength)
	_, err := r.Body.Read(buf)
	if err != nil && err != io.EOF {
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("responce marshal error: %s", err)
	}
	return
}
