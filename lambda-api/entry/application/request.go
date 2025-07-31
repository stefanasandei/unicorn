package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/entry/repository/worker"
	"github.com/redis/go-redis/v9"
)

type FailedMessage struct {
	Status  string `json:"status"`
	Message string `json:"output"`
}

func FailIfError(err error, w http.ResponseWriter, msg string) bool {
	if err == nil {
		return false
	}

	sMsg, err := json.Marshal(FailedMessage{
		Status:  "failed",
		Message: fmt.Sprintf("%s: %s", msg, err),
	})

	w.WriteHeader(http.StatusServiceUnavailable)

	_, err = w.Write(sMsg)
	if err != nil {
		log.Printf("Failed to write a response: %s", err)
	}

	return true
}

func SetupStreamingResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/stream+json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	return nil
}

func ChooseWorker(rdb *redis.Client, c context.Context) (string, error) {
	// get all workers from redis
	workersRepo := worker.RedisWorkerRepository{
		Client: rdb,
	}

	workers, err := workersRepo.QueryWorkers(c)
	if err != nil {
		return "", err
	}

	// choose the best one, the one with the least work at the moment
	workerId, err := workersRepo.GetBestWorker(workers)
	if err != nil {
		return "", err
	}

	return workerId, nil
}
