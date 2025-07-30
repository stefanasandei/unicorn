package worker

import (
	"context"

	"binomeway.com/worker/model"
	"github.com/google/uuid"
)

type WorkerRepository interface {
	RegisterWorker(ctx context.Context, id uuid.UUID) model.Worker

	Heartbeat(ctx context.Context, worker *model.Worker)

	UnregisterWorker(ctx context.Context, worker model.Worker)
}
