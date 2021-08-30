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
	SearchAvailabilityByDates(start, end time.Time, roomID int) (bool, error)
}
