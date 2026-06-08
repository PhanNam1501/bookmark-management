package validator

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PhanNam1501/bookmark-management/internal/service/queue"
)

const (
	MaxDescriptionLen = 500
	MaxURLLen         = 2048
)

// ValidateImportData validates CSV import data
func ValidateImportData(data []*queue.ImportBookmarkInput) error {
	if len(data) == 0 {
		return fmt.Errorf("CSV file is empty")
	}

	if len(data) > 1000 {
		return fmt.Errorf("too many bookmarks, max 1000 per import")
	}

	seenURLs := make(map[string]bool)

	for i, bookmark := range data {
		if err := validateBookmark(bookmark, i); err != nil {
			return err
		}

		// Check duplicate
		if seenURLs[bookmark.Url] {
			return fmt.Errorf("duplicate URL at row %d: %s", i+2, bookmark.Url)
		}
		seenURLs[bookmark.Url] = true
	}

	return nil
}

func validateBookmark(b *queue.ImportBookmarkInput, rowIndex int) error {
	rowNum := rowIndex + 2 // CSV headers = row 1

	// Validate URL
	if strings.TrimSpace(b.Url) == "" {
		return fmt.Errorf("row %d: URL is required", rowNum)
	}

	if len(b.Url) > MaxURLLen {
		return fmt.Errorf("row %d: URL is too long (max %d characters)", rowNum, MaxURLLen)
	}

	if _, err := url.ParseRequestURI(b.Url); err != nil {
		return fmt.Errorf("row %d: invalid URL format", rowNum)
	}

	// Validate Description
	if strings.TrimSpace(b.Description) == "" {
		return fmt.Errorf("row %d: Description is required", rowNum)
	}

	if len(b.Description) > MaxDescriptionLen {
		return fmt.Errorf("row %d: Description is too long (max %d characters)", rowNum, MaxDescriptionLen)
	}

	return nil
}
