package main

import (
	"log"
	"net/http"

	"lukasaebi.ch/cinema-booking-system/internal/adapters/redis"
	"lukasaebi.ch/cinema-booking-system/internal/booking"
	"lukasaebi.ch/cinema-booking-system/internal/utils"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /movies", listMovies)

	store := booking.NewRedisStore(redis.NewClient("localhost:6379"))
	svc := booking.NewService(store)
	bookingHandler := booking.NewHandler(svc)

	mux.HandleFunc("GET /movies/{movieID}/bookings", bookingHandler.ListBookedSeats)

	mux.Handle("GET /", http.FileServer(http.Dir("static")))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

type movieResponse struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Rows         int    `json:"rows"`
	SeathsPerRow int    `json:"seats_per_row"`
}

var movies = []movieResponse{
	{ID: "inception", Title: "Inception", Rows: 5, SeathsPerRow: 8},
	{ID: "dune", Title: "Dune: Part Two", Rows: 4, SeathsPerRow: 6},
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, movies)
}
