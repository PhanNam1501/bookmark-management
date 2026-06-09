package worker

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PhanNam1501/bookmark-management/internal/repository/queue"
	"github.com/rs/zerolog/log"
)

type Engine interface {
	Start(ctx context.Context)
}

type Handler interface {
	Handle(ctx context.Context, message []byte) error
}

type Queue interface {
	PopMessage(ctx context.Context) ([]byte, error)
}

type engine struct {
	// queue
	queue Queue
	// handler
	handler Handler
	run     bool
	sigChan chan os.Signal
}

func NewEngine(queue Queue, handler Handler) Engine {
	return &engine{
		queue:   queue,
		handler: handler,
		run:     false,
		sigChan: make(chan os.Signal, 1),
	}
}

const workerSleepDuration = 1 * time.Second

func (e *engine) Start(ctx context.Context) {
	println("starting worker...")
	workerPool := NewPool(ctx, e.handler, 4)
	signal.Notify(e.sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	e.run = true
	for e.run {
		select {
		case sig := <-e.sigChan:
			log.Info().Msgf("received signal: %s. Shutting down worker Engine", sig)
			e.run = false
		default:
			msg, err := e.queue.PopMessage(ctx)
			if err != nil {
				if errors.Is(err, queue.ErrorNoMessage) {
					time.Sleep(workerSleepDuration)
					continue
				}

				log.Error().Err(err).Msg("failed to pop message")
				continue
			}
			workerPool.Consume(msg)
		}
	}
	workerPool.Close()
}
