package app

import (
	"redisQueue/internal/app/rest"
	"redisQueue/storage/redis"
)

type App struct {
	RestApp rest.RestApp
}

func New(port string, storageUrl string, key string, maxLen int) *App {

	redisProvider := storage.New(storageUrl, key, maxLen)

	postApp := rest.New(port, redisProvider)
	return &App{
		RestApp: *postApp,
	}
}
