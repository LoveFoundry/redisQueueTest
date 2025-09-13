package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"redisQueue/domains/models"
	"redisQueue/storage/redis"

	"github.com/gorilla/mux"
)

type ServerApi struct {
	provider RedisProvider
	port     string
}
type RedisProvider interface {
	Add(ctx context.Context, repeatNum int, m string) error
}

func New(provider *storage.RedisProvider, port string) *ServerApi {
	api := ServerApi{
		provider: provider,
		port:     port,
	}
	return &api
}

func (s *ServerApi) ConfigureRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/addMessage", s.AddMessage).Methods("POST")

	return r
}
func (s *ServerApi) AddMessage(w http.ResponseWriter, r *http.Request) {
	var req models.Msg
	slog.Info("Received request for add message")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	if err := addMessageValidation(req, w); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(req)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.provider.Add(ctx, req.RepeatNum, string(b))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}

func addMessageValidation(req models.Msg, w http.ResponseWriter) error {
	if req.Message == "" {
		http.Error(w, "Message обязателен", http.StatusBadRequest)
		return fmt.Errorf("Message обязателен")
	}
	if req.RepeatNum == 0 {
		http.Error(w, "repeat num должен быть или 1 или 7", http.StatusBadRequest)
		return fmt.Errorf("repeat num должен быть или 1 или 7")
	}
	if req.Message == "" {
		http.Error(w, "userId обязателен", http.StatusBadRequest)
		return fmt.Errorf("userId обязателен")
	}
	return nil
}
