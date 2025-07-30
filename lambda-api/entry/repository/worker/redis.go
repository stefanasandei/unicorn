package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisWorkerRepository struct {
	Client *redis.Client
}

type SimplifiedWorker struct {
	ID       uuid.UUID `json:"id"`
	CPUUsage float64   `json:"cpu_usage"`
}

func (rdb *RedisWorkerRepository) QueryWorkers(ctx context.Context) ([]SimplifiedWorker, error) {
	var workers []SimplifiedWorker
	worker := SimplifiedWorker{}

	// get a list of all worker IDs
	ids, err := rdb.Client.SMembers(ctx, "workers").Result()
	if err != nil {
		return nil, err
	}

	// now query each worker
	for _, id := range ids {
		workerId := fmt.Sprintf("worker:%s", id)

		value, err := rdb.Client.Get(ctx, workerId).Result()
		if err != nil {
			return nil, err
		}

		// unmarshal the value
		err = json.Unmarshal([]byte(value), &worker)
		if err != nil {
			return nil, err
		}

		workers = append(workers, worker)
	}

	return workers, nil
}

func (rdb *RedisWorkerRepository) GetBestWorker(workers []SimplifiedWorker) (string, error) {
	if len(workers) == 0 {
		return "", fmt.Errorf("empty workers array")
	}

	bestWorker := SimplifiedWorker{
		CPUUsage: math.MaxFloat64,
	}

	for _, worker := range workers {
		if worker.CPUUsage < bestWorker.CPUUsage {
			bestWorker = worker
		}
	}

	return bestWorker.ID.String(), nil
}
