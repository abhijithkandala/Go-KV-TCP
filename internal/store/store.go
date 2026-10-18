package store

import "errors"

var ErrEntryNotFound = errors.New("entry not found")
var ErrUnknownOperation = errors.New("unknown operation")

type Store interface {
	Get(key string) (string, error)
	Set(key, value string)
	Del(key string) error
}
