package rest

import (
	"log"
	"net/http"
	"os"
	"redisQueue/internal/server/rest"
	storage "redisQueue/storage/redis"

	"github.com/rs/cors"
)

type RestApp struct {
	server *rest.ServerApi
	port   string
}

func New(port string, provider *storage.RedisProvider) *RestApp {

	server := rest.New(provider, port)

	return &RestApp{
		server: server,
		port:   port,
	}
}

func (app *RestApp) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *RestApp) Run() error {
	log.Printf("asdfasdf %s", app.port)
	r := app.server.ConfigureRoutes()
	http.Handle("/", cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}).Handler(r))

	go func() {
		if err := http.ListenAndServe(app.port, nil); err != nil {
			log.Printf("Rest rest listening on port %s", app.port)
		}
	}()
	return nil
}
func (app *RestApp) Stop() {
	log.Printf("Rest rest shutting down")

	os.Exit(0)
}
