package queue

import (
	"context"

	"github.com/PhanNam1501/bookmark-management/internal/repository/queue"
)

type Service interface {
	SendImportBookmarkJob(ctx context.Context, uid string, bookmarkInput []*ImportBookmarkInput) error
}

type queueService struct {
	q queue.Repository
}

func NewService(q queue.Repository) Service {
	return &queueService{
		q: q,
	}
}
