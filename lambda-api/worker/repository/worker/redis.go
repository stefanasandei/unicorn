package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"

	"binomeway.com/worker/model"
	"github.com/redis/go-redis/v9"
)

type RedisWorkerRepository struct {
	Client *redis.Client
}

func (rdb *RedisWorkerRepository) RegisterWorker(ctx context.Context, id uuid.UUID, runtimesDir string) model.Worker {
	worker := model.NewWorker(id, runtimesDir)

	// encode the worker into json to be later serialized into a redis hashmap
	data, err := json.Marshal(worker)
	if err != nil {
		log.Panicf("Failed to encode order: %s", err)
	}

	txn := rdb.Client.TxPipeline()

	// create a hashmap of the worker's data
	res := txn.SetNX(ctx, fmt.Sprintf("worker:%s", id.String()), string(data), 0)
	if err := res.Err(); err != nil {
		txn.Discard()
		log.Panicf("Failed to set: %s", err)
	}

	// add the uuid of the worker to the list of workers
	if err := txn.SAdd(ctx, "workers", id.String()).Err(); err != nil {
		txn.Discard()
		log.Panicf("Failed to add to workers set: %s", err)
	}

	// execute the pipeline
	if _, err := txn.Exec(ctx); err != nil {
		log.Panicf("Failed to exec the transaction: %s", err)
	}

	return worker
}

func (rdb *RedisWorkerRepository) Heartbeat(ctx context.Context, worker *model.Worker) {
	worker.Update()

	data, err := json.Marshal(worker)
	if err != nil {
		log.Panicf("Failed to encode order: %s", err)
	}

	txn := rdb.Client.TxPipeline()

	res := txn.Set(ctx, fmt.Sprintf("worker:%s", worker.ID.String()), string(data), 0)
	if err := res.Err(); err != nil {
		txn.Discard()
		log.Panicf("Failed to set: %s", err)
	}

	// execute the pipeline
	if _, err := txn.Exec(ctx); err != nil {
		log.Panicf("Failed to exec the transaction: %s", err)
	}
}

func (rdb *RedisWorkerRepository) UnregisterWorker(ctx context.Context, worker model.Worker) {
	txn := rdb.Client.TxPipeline()

	// delete the worker from Redis
	res := txn.Del(ctx, fmt.Sprintf("worker:%s", worker.ID.String()))
	if err := res.Err(); err != nil {
		txn.Discard()
		log.Panicf("Failed to delete: %s", err)
	}

	// add the uuid of the worker to the list of workers
	if err := txn.SRem(ctx, "workers", worker.ID.String()).Err(); err != nil {
		txn.Discard()
		log.Panicf("Failed to remove from workers set: %s", err)
	}

	// execute the pipeline
	if _, err := txn.Exec(ctx); err != nil {
		log.Panicf("Failed to exec the transaction: %s", err)
	}
}
