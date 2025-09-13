package main

import (
	"log"
	"os"
	"os/signal"
	"redisQueue/internal/app"
	"strconv"
	"syscall"
)

func main() {
	var (
		port         = os.Getenv("PORT")
		redisHost    = os.Getenv("REDIS_ADDR")
		key          = os.Getenv("QUEUE_KEY")
		maxLenString = os.Getenv("MAXLEN")
	)
	maxLen, err := strconv.Atoi(maxLenString)
	if err != nil {
		log.Fatal("Max length invalid: ", err)
	}
	application := app.New(port, redisHost, key, maxLen)

	application.RestApp.MustRun()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sign := <-stop

	log.Printf("Signal", sign.String())

	application.RestApp.Stop()

	log.Printf("Shutting down")
}
