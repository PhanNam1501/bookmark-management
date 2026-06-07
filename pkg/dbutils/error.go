package dbutils

import (
	"errors"
	"strings"
)

var errorFilter = []func(error) (bool, error){
	filterRecordNotFound,
	filterDuplicateError,
	filterForeignKeyError,
}

func CatchDBError(err error) error {
	if err == nil {
		return nil
	}
	for _, filter := range errorFilter {
		match, filterErr := filter(err)
		if match {
			return filterErr
		}
	}

	return err
}

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicate      = errors.New("duplicate")
	ErrForeignKey     = errors.New("foreign key")
)

func filterRecordNotFound(err error) (bool, error) {
	return strings.Contains(strings.ToLower(err.Error()), "record not found"), ErrRecordNotFound
}

func filterDuplicateError(err error) (bool, error) {
	return strings.Contains(strings.ToLower(err.Error()), "duplicate"), ErrDuplicate
}

func filterForeignKeyError(err error) (bool, error) {
	return strings.Contains(strings.ToLower(err.Error()), "foreign key constraint"), ErrForeignKey
}
