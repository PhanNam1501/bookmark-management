package workerhandler

import (
	"context"
	"time"
)

type TestHandler struct {
}

func (t *TestHandler) Handle(ctx context.Context, message []byte) error {
	time.Sleep(1 * time.Second)
	println(string(message))
	return nil
}
