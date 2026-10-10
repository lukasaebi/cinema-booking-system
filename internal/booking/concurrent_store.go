package booking

import "sync"

type ConcurrentStore struct {
	bookings map[string]Booking // seats -> booking
	sync.RWMutex
}

func NewConcurrentStore() *ConcurrentStore {
	return &ConcurrentStore{
		bookings: map[string]Booking{},
	}
}

func (s *ConcurrentStore) Book(b Booking) error {
	s.Lock() // blocks everyone
	defer s.Unlock()

	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}
	s.bookings[b.SeatID] = b
	return nil
}

func (s *ConcurrentStore) ListBookings(movieID string) []Booking {
	s.RLock() // blocks writers, allows other RLock()ers
	defer s.RUnlock()

	var movies []Booking
	for _, b := range s.bookings {
		if b.MovieID == movieID {
			movies = append(movies, b)
		}
	}
	return movies
}
