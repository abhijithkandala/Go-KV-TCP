package store

type MutexStore struct {
	data map[string]string
}

func NewMutexStore() *MutexStore {
	return &MutexStore{make(map[string]string)}
}

func (s *MutexStore) Set(key, value string) {
	s.data[key] = value
	return
}

func (s *MutexStore) Get(key string) (string, error) {
	val, isPresent := s.data[key]
	if !isPresent {
		return "", ErrEntryNotFound
	}
	return val, nil
}
func (s *MutexStore) Del(key string) error {
	_, isPresent := s.data[key]
	if !isPresent {
		return ErrEntryNotFound
	}
	delete(s.data, key)
	return nil
}

var _ Store = (*MutexStore)(nil)
