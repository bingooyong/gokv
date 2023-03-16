package server

import (
	"encoding/json"
	"errors"
)

func (s *Server) Get(key string) (interface{}, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if v, ok := s.storage[key]; ok {
		return v, nil
	}

	return nil, errors.New("not found")
}

func (s *Server) Set(key string, val interface{}) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.storage[key] = val

	b, err := json.Marshal([]*Update{
		{
			Action: "set",
			Data: map[string]interface{}{
				key: val,
			},
		},
	})

	if err != nil {
		return err
	}

	s.broadcasts.QueueBroadcast(&broadcast{
		msg:    append([]byte("d"), b...),
		notify: nil,
	})

	return nil
}

func (s *Server) Delete(key string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, ok := s.storage[key]; ok {
		delete(s.storage, key)
	}

	b, err := json.Marshal([]*Update{{
		Action: "del",
		Data: map[string]interface{}{
			key: nil,
		},
	}})

	if err != nil {
		return err
	}

	s.broadcasts.QueueBroadcast(&broadcast{
		msg:    append([]byte("d"), b...),
		notify: nil,
	})

	return nil
}
