package internalhttp

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const defaultInterval = 2 * time.Second

type Response struct {
	Data  interface{} `json:"data"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type Service struct {
	sync.RWMutex
	Stats    map[uint32]uint32
	Interval time.Duration
}

func NewService() *Service {
	return &Service{
		Stats:    make(map[uint32]uint32),
		Interval: defaultInterval,
	}
}

func (s *Service) Index(w http.ResponseWriter, r *http.Request) {
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
	return
}

func WriteResponse(w http.ResponseWriter, resp *Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("responce marshal error: %s", err)
	}

	return
}
