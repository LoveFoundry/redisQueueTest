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
	Get(ctx context.Context) ([]models.Msg, error)
}

func New(provider *storage.RedisProvider, port string) *ServerApi {
	api := ServerApi{
		provider: provider,
		port:     port,
	}
	return &api
}

func (s *ServerApi) ConfigureRoutes() *mux.Router {
	slog.Info("routes initialized")

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/addMessage", s.AddMessage).Methods("POST")
	api.HandleFunc("/getMessages", s.GetMessages).Methods("GET")

	return r
}
func (s *ServerApi) AddMessage(w http.ResponseWriter, r *http.Request) {
	var req models.MsgRequest
	slog.Info("Received request for add message")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	ctx := r.Context()
	if err := addMessageValidation(req, w); err != nil {
		slog.Info(err.Error())
		return
	}

	b, err := json.Marshal(models.Msg{UserID: req.UserID, Message: req.Message})

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

func (s *ServerApi) GetMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	msg, err := s.provider.Get(ctx)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"messages": msg,
	})
	if err != nil {
		return
	}
}

func addMessageValidation(req models.MsgRequest, w http.ResponseWriter) error {
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
