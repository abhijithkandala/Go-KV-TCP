package store

import (
	"errors"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	s := NewMutexStore()
	s.Set("foo", "bar")

	got, err := s.Get("foo")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got != "bar" {
		t.Errorf("got %q, want %q", got, "bar")
	}
}

func TestGetMissingKey(t *testing.T) {
	s := NewMutexStore()

	_, err := s.Get("nope")
	if !errors.Is(err, ErrEntryNotFound) {
		t.Errorf("got %v, want ErrEntryNotFound", err)
	}
}

func TestSetOverwrites(t *testing.T) {
	s := NewMutexStore()
	s.Set("foo", "bar")
	s.Set("foo", "baz")

	got, err := s.Get("foo")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got != "baz" {
		t.Errorf("got %q, want %q", got, "baz")
	}
}

func TestDelExistingKey(t *testing.T) {
	s := NewMutexStore()
	s.Set("foo", "bar")

	if err := s.Del("foo"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if _, err := s.Get("foo"); !errors.Is(err, ErrEntryNotFound) {
		t.Errorf("got %v, want ErrEntryNotFound after delete", err)
	}
}

func TestDelMissingKey(t *testing.T) {
	s := NewMutexStore()

	if err := s.Del("nope"); !errors.Is(err, ErrEntryNotFound) {
		t.Errorf("got %v, want ErrEntryNotFound", err)
	}
}

func TestEmptyValueIsNotMissing(t *testing.T) {
	s := NewMutexStore()
	s.Set("k", "")

	got, err := s.Get("k")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("got %q, want empty string", got)
	}
}
