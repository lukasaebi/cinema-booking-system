package booking

import (
	"net/http"

	"lukasaebi.ch/cinema-booking-system/internal/utils"
)

type handler struct {
	svc *Service
}

func NewHandler(svc *Service) *handler {
	return &handler{svc}
}

func (h *handler) ListBookedSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	bookings := h.svc.ListBookings(movieID)

	seats := make([]SeatInfo, 0, len(bookings))
	for _, b := range bookings {
		seats = append(seats, SeatInfo{
			SeatID: b.SeatID,
			UserID: b.UserID,
			Booked: true,
		})
	}
	utils.WriteJSON(w, http.StatusOK, seats)
}

type SeatInfo struct {
	SeatID string `json:"seat_id"`
	UserID string `json:"user_id"`
	Booked bool   `json:"booked"`
}
