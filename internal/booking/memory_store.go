package booking

type MemoryStore struct {
	bookings map[string]Booking // seats -> booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookings: map[string]Booking{},
	}
}

func (s *MemoryStore) Book(b Booking) error {
	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}
	s.bookings[b.SeatID] = b
	return nil
}

func (s *MemoryStore) ListBookings(movieID string) []Booking {
	var movies []Booking
	for _, b := range s.bookings {
		if b.MovieID == movieID {
			movies = append(movies, b)
		}
	}
	return movies
}
