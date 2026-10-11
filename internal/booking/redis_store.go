package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const defaultHoldTTL = time.Minute * 2

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func (s *RedisStore) Book(b Booking) error {
	session, err := s.hold(b)
	if err != nil {
		return err
	}

	log.Printf("Session booked %v", session)

	return nil
}

func (s *RedisStore) ListBookings(movieID string) []Booking {
	return []Booking{}
}

func (s *RedisStore) hold(b Booking) (Booking, error) {
	ctx := context.Background()
	key := fmt.Sprintf("seat:%s:%s", b.MovieID, b.SeatID)

	b.ID = uuid.New().String()
	b.Status = "held"
	b.ExpiresAt = time.Now().Add(defaultHoldTTL)
	val, err := json.Marshal(b)
	if err != nil {
		return Booking{}, err
	}

	// entry for a seat
	ok, err := s.rdb.SetNX(ctx, key, val, defaultHoldTTL).Result()
	if err != nil {
		return Booking{}, err
	}
	if !ok {
		return Booking{}, ErrSeatAlreadyBooked
	}

	// entry for a session -> is our index for looking up seat entries
	err = s.rdb.Set(ctx, sessionKey(b.ID), val, defaultHoldTTL).Err()
	if err != nil {
		s.rdb.Del(ctx, key) // roll back the seat lock
		return Booking{}, err
	}

	return b, nil
}
