package db

const (
	DBURI     = "mongodb://localhost:27017"
	DBNAME    = "hotel-reservation"
	TESTDBURI = "mongodb://localhost:27017"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
