package modals

import (
	"time"
)

//User is the user modal
type User struct {
	ID int
	FirstName string
	LastName string
	Email string
	Password string
	AccessLevel int
	CreatedAt time.Time
	UpdatedAt time.Time
}

//Room is the room modal
type Room struct{
	ID int
	RoomName string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//Restriction is the Restriction modal
type Restriction struct{
	ID int
	RestrictionName string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//Reservation is the Reservation modal
type Reservation struct{
	ID int
	FirstName string
	LastName string
	Email string
	Phone string
	StartDate time.Time
	EndDate time.Time
	RoomID int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room Room
}

//RoomRestriction is the RoomRestriction modal
type RoomRestriction struct{
	ID int
	StartDate time.Time
	EndDate time.Time
	RoomID int
	ReservationID int
	RestrictionID int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room Room
	Reservation Reservation
	Restriction Restriction
}
