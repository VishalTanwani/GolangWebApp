package repository

import (
	"github.com/VishalTanwani/GolangWebApp/internal/modals"
	"time"
)

//DatabaseRepo interface will hold all db functions
type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res modals.Reservation) (int, error)
	InsertRoomRestriction(res modals.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailbilityForAllRooms(start, end time.Time) ([]modals.Room, error)
	GetRoomByID(id int) (modals.Room, error)
	GetUserByID(id int) (modals.User, error)
	UpdateUser(user modals.User) error
	Authenticate(email, password string) (int, string, error)
	GetAllReservations() ([]modals.Reservation, error)
	GetAllNewReservations() ([]modals.Reservation, error)
	GetReservationByID(id int) (modals.Reservation, error)
}
