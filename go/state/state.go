package state

import "sync"

var (
	GraphToken         = State[string]{}
	AppSecret          = State[string]{}
	PhoneNumber        = State[string]{}
	PhoneNumberID      = State[string]{}
	WebhookURL         = State[string]{}
	WebhookVerifyToken = State[string]{}
)

type State[T any] struct {
	value T
	lock  sync.Mutex
}

func (s *State[T]) Set(value T) {
	s.lock.Lock()
	s.value = value
	s.lock.Unlock()
}

func (s *State[T]) Get() T {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.value
}
