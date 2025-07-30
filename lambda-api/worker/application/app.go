package application

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"

	"binomeway.com/common/broker"
	"binomeway.com/worker/model"
	"binomeway.com/worker/repository/worker"
)

type App struct {
	config Config
	rdb    *redis.Client
	broker broker.MessageBroker
	worker model.Worker
}

func New(config Config) *App {
	app := &App{
		config: config,
		rdb: redis.NewClient(&redis.Options{
			Addr: config.RedisAddress,
		}),
		broker: &broker.RabbitMQMessageBroker{},
	}

	return app
}

func (app *App) Start(ctx context.Context) error {
	err := app.broker.Connect(app.config.RabbitMQAddress)
	if err != nil {
		return err
	}

	err = app.broker.CreateQueue(app.config.ID.String())
	if err != nil {
		return err
	}

	msgs, err := app.broker.Consume(app.config.ID.String())
	if err != nil {
		return err
	}

	// start the redis connection
	err = app.setupRedis(ctx)
	if err != nil {
		return err
	}

	// close the Redis connection when we're done
	defer func() {
		err := app.closeRedis()
		if err != nil {
			return
		}
	}()

	app.setupCronJobs()

	log.Printf("Started worker: %s", app.config.ID.String())
	log.Printf("[*] Waiting for messages. To exit press CTRL+C")

	// create multiple threads to handle messages
	threadsNum := 10

	// start a communication channel for the error
	var wg sync.WaitGroup
	ch := make(chan error, threadsNum)

	// main loop
	for i := 0; i < threadsNum; i++ {
		wg.Add(1)

		log.Printf("Thread %d started.", i)

		go func(workerID int) {
			defer wg.Done()

			for d := range msgs {
				err := app.HandleQueueMessage(d)
				if err != nil {
					log.Printf("Worker %d: Failed to handle a message: %s", workerID, err)
					ch <- err
				}
			}
		}(i)
	}

	// graceful shutdown
	select {
	case err := <-ch:
		return err

	case <-ctx.Done():
		close(ch)

		_, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return app.broker.Close()
	}
}

func (app *App) setupRedis(ctx context.Context) error {
	err := app.rdb.Ping(ctx).Err()

	if err != nil {
		log.Panicf("Failed to connect to redis: %s", err)
		return err
	}

	repo := worker.RedisWorkerRepository{
		Client: app.rdb,
	}

	app.worker = repo.RegisterWorker(ctx, app.config.ID, app.config.RuntimesDir)

	return nil
}

func (app *App) closeRedis() error {
	repo := worker.RedisWorkerRepository{
		Client: app.rdb,
	}

	repo.UnregisterWorker(context.Background(), app.worker)

	if err := app.rdb.Close(); err != nil {
		log.Panicf("Failed to close redis: %s", err)
		return err
	}

	log.Printf("Closed the Redis connection.")

	return nil
}

func (app *App) setupCronJobs() {
	repo := worker.RedisWorkerRepository{
		Client: app.rdb,
	}

	heartbeatInterval := 30
	timeInterval := fmt.Sprintf("@every %ds", heartbeatInterval)

	c := cron.New()

	_, err := c.AddFunc(timeInterval, func() {
		repo.Heartbeat(context.Background(), &app.worker)
	})
	if err != nil {
		log.Printf("Failed to create the heardbeat function.")
		return
	}

	c.Start()
}
