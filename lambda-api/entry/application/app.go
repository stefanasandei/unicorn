package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	common "github.com/common/model"

	"github.com/common/broker"
	"github.com/entry/model"
	"github.com/redis/go-redis/v9"
)

type ExecutionContext struct {
	incomingMessages []common.WorkerResponseWrapper
	backgroundCtx    context.Context
	messageMutex     sync.Mutex
}

type App struct {
	router http.Handler

	redisDB *redis.Client
	broker  broker.MessageBroker
	reply   *model.ReplyQueue

	ctx ExecutionContext

	config Config
}

func New(config Config) *App {
	app := &App{
		redisDB: redis.NewClient(&redis.Options{
			Addr: config.RedisAddress,
		}),
		broker: &broker.RabbitMQMessageBroker{},
		ctx: ExecutionContext{
			incomingMessages: make([]common.WorkerResponseWrapper, 0),
			backgroundCtx:    context.Background(),
		},
		config: config,
	}

	app.loadRoutes()

	return app
}

func (a *App) Start(ctx context.Context) error {
	err := a.broker.Connect(a.config.RabbitMQAddress)
	if err != nil {
		return err
	}

	a.reply, err = model.NewReplyQueue(a.broker)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:    a.config.ServerAddr,
		Handler: a.router,
	}

	err = a.redisDB.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	defer func() {
		if err := a.redisDB.Close(); err != nil {
			fmt.Println("failed to close redis", err)
		}
		fmt.Println("closed redis!")
	}()

	log.Printf("Starting the code execution service on http://%s\n", a.config.ServerAddr)

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := a.reply.Close()
		if err != nil {
			return err
		}

		return server.Shutdown(timeout)
	}
}
