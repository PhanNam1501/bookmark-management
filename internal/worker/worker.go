package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Pool struct {
	handler     Handler
	numWorkers  int
	messageChan chan []byte
	errChan     chan *worker
	wg          *sync.WaitGroup
}

func NewPool(ctx context.Context, handler Handler, numWorkers int) *Pool {
	messageChan := make(chan []byte, numWorkers)
	errChan := make(chan *worker, numWorkers)
	pool := &Pool{
		handler:     handler,
		numWorkers:  numWorkers,
		messageChan: messageChan,
		errChan:     errChan,
	}
	pool.init(ctx)

	fmt.Printf("starting worker pool with %d workers", numWorkers)
	return pool
}

func (p *Pool) Consume(message []byte) {
	p.messageChan <- message
}

func (p *Pool) init(ctx context.Context) {
	for i := 1; i <= p.numWorkers; i++ {
		worker := &worker{
			id:          i,
			handler:     p.handler,
			messageChan: p.messageChan,
			errChan:     p.errChan,
			wg:          p.wg,
		}
		worker.wg.Add(1)

		go worker.Work(ctx)
	}

	go func() {
		for w := range p.errChan {
			log.Error().Msgf("worker %d encounter a panic", w.id)
			time.Sleep(2 * time.Second)
			log.Info().Msgf("worker %d restarting", w.id)
			go w.Work(ctx)
		}
	}()
}

type worker struct {
	id          int
	handler     Handler
	messageChan <-chan []byte
	err         error
	errChan     chan<- *worker
	wg          *sync.WaitGroup
}

func (w *worker) Work(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				w.err = err
			} else {
				w.err = fmt.Errorf("worker panicked: %v", r)
			}
			w.errChan <- w
		} else {
			log.Info().Msgf("worker %d exits normally", w.id)
			w.wg.Done()
		}
	}()

	for {
		msg, ok := <-w.messageChan
		if !ok {
			log.Info().Msgf("worker %d is closing", w.id)
			return
		}
		log.Info().Msgf("worker %d received message and processing it", w.id)

		err := w.handler.Handle(ctx, msg)
		if err != nil {
			log.Error().Err(err).Msgf("worker %d failed to handle message", w.id)
		} else {
			log.Info().Msgf("worker %d handled message", w.id)
		}
	}
}

func (p *Pool) Close() {
	close(p.messageChan)
	p.wg.Wait()
	close(p.errChan)
}
