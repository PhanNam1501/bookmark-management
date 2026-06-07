package bookmark

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/PhanNam1501/bookmark-management/internal/model"
	"github.com/PhanNam1501/bookmark-management/internal/repository/cache"
	"github.com/rs/zerolog/log"
)

const (
	getBookmarkCacheGroupFormat = "get_bookmarks_user_%s"
	getBookmarkCacheKeyFormat   = "page_%d_limit_%d"
	getBookmarkCacheDuration    = 24 * time.Hour
)

type bookmarkCacheService struct {
	s     Service
	cache cache.Repo
}

func NewServiceWithCache(s Service, cache cache.Repo) Service {
	return &bookmarkCacheService{
		s:     s,
		cache: cache,
	}
}

func (s *bookmarkCacheService) GetBookmarks(ctx context.Context, userID string, limit, page int) (*GetBookmarksResult, error) {
	groupKey := fmt.Sprintf(getBookmarkCacheGroupFormat, userID)
	cacheKey := fmt.Sprintf(getBookmarkCacheKeyFormat, page, limit)

	// check cache
	cacheData, err := s.cache.GetCacheData(ctx, groupKey, cacheKey)
	if err == nil && len(cacheData) > 0 {
		result := &GetBookmarksResult{}
		err := json.Unmarshal(cacheData, result)
		if err == nil {
			return result, nil
		}
	}

	result, err := s.s.GetBookmarks(ctx, userID, limit, page)
	if err != nil {
		return nil, err
	}

	resultBytes, err := json.Marshal(result)
	if err == nil {
		err = s.cache.SetCacheData(ctx, groupKey, cacheKey, resultBytes, getBookmarkCacheDuration)
		if err != nil {
			log.Warn().Err(err).Msg("Can't cache to redis")
		}
	}

	return result, nil
}

func (s *bookmarkCacheService) CreateBookmark(ctx context.Context, userid, url, description string) (*model.Bookmark, error) {
	err := s.cache.DeleteCacheData(ctx, fmt.Sprintf(getBookmarkCacheGroupFormat, userid))
	if err != nil {
		return nil, err
	}

	return s.s.CreateBookmark(ctx, userid, url, description)
}

func (s *bookmarkCacheService) UpdateBookmark(ctx context.Context, bookmarkID, userID string, url, description string) error {
	err := s.cache.DeleteCacheData(ctx, fmt.Sprintf(getBookmarkCacheGroupFormat, userID))
	if err != nil {
		return err
	}

	return s.s.UpdateBookmark(ctx, bookmarkID, userID, url, description)
}

func (s *bookmarkCacheService) DeleteBookmark(ctx context.Context, bookmarkID, userID string) error {
	err := s.cache.DeleteCacheData(ctx, fmt.Sprintf(getBookmarkCacheGroupFormat, userID))
	if err != nil {
		return err
	}

	return s.s.DeleteBookmark(ctx, bookmarkID, userID)
}
