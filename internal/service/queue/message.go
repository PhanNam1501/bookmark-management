package queue

import (
	"context"
	"encoding/json"

	"github.com/PhanNam1501/bookmark-management/pkg/array"
)

const BatchSize = 20

type ImportMessage struct {
	UID       string                 `json:"uid"`
	Bookmarks []*ImportBookmarkInput `json:"bookmarks"`
}

type ImportBookmarkInput struct {
	Description string `csv:"description" validate:"lte=255" json:"description"`
	Url         string `csv:"url" validate:"required,url,lte=2048" json:"url"`
}

func (s *queueService) SendImportBookmarkJob(ctx context.Context, uid string, bookmarkInput []*ImportBookmarkInput) error {
	batches := array.SplitIntoBatches[*ImportBookmarkInput](bookmarkInput, BatchSize)

	for _, batch := range batches {
		err := s.sendJob(ctx, uid, batch)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *queueService) sendJob(ctx context.Context, uid string, bookmarkInput []*ImportBookmarkInput) error {
	message := ImportMessage{
		UID:       uid,
		Bookmarks: bookmarkInput,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return s.q.PushMessage(ctx, messageBytes)
}
