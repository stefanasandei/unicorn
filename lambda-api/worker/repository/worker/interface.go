package worker

import (
	"context"

	"github.com/google/uuid"
	"github.com/worker/model"
)

type WorkerRepository interface {
	RegisterWorker(ctx context.Context, id uuid.UUID) model.Worker

	Heartbeat(ctx context.Context, worker *model.Worker)

	UnregisterWorker(ctx context.Context, worker model.Worker)
}
