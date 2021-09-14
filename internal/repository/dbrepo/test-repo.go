package dbrepo

import (
	"errors"
	"github.com/VishalTanwani/GolangWebApp/internal/modals"
	"time"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

//InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res modals.Reservation) (int, error) {
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

//InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(res modals.RoomRestriction) error {
	if res.RoomID == 3 {
		return errors.New("some error")
	}
	return nil
}

//SearchAvailbilityByDates return status of availability of room
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	if roomID == 2 {
		return false, nil
	}
	if roomID == 3 {
		return false, errors.New("some error")
	}
	return true, nil
}

//SearchAvailbilityForAllRooms will give all the rooms available in given dates
func (m *testDBRepo) SearchAvailbilityForAllRooms(start, end time.Time) ([]modals.Room, error) {
	var rooms []modals.Room
	sd := start.Format("2006-01-02")
	if "2021-11-01" == sd {
		return []modals.Room{{ID: 1, RoomName: "masdbf"}}, nil
	}
	return rooms, nil
}

//GetRoomByID will get room by id
func (m *testDBRepo) GetRoomByID(id int) (modals.Room, error) {
	var room modals.Room
	if id > 2 {
		return room, errors.New("some error")
	}
	return room, nil
}

//GetUserByID will give user object
func (m *testDBRepo) GetUserByID(id int) (modals.User, error) {
	var user modals.User

	return user, nil
}

//UpdateUser will update the given user
func (m *testDBRepo) UpdateUser(user modals.User) error {
	return nil
}

//Authenticate it will Authenticate a user
func (m *testDBRepo) Authenticate(email, password string) (int, string, error) {
	if email == "admin@admin.com" {
		return 1, "", nil
	}
	return 1,"",errors.New("invalid credentails")

}

//GetAllReservations will return all reservation
func (m *testDBRepo) GetAllReservations() ([]modals.Reservation, error) {
	var reservations []modals.Reservation
	return reservations, nil
}

//GetAllNewReservations will return all new reservation
func (m *testDBRepo) GetAllNewReservations() ([]modals.Reservation, error) {
	var reservations []modals.Reservation
	return reservations, nil
}

//GetReservationByID will return one reservation by id
func (m *testDBRepo) GetReservationByID(id int) (modals.Reservation, error) {
	var reservation modals.Reservation
	return reservation, nil
}

//UpdateReservation will update the given reservation
func (m *testDBRepo) UpdateReservation(reservation modals.Reservation) error {
	return nil
}

//DeleteReservation will delete the reservation
func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

//UpdateProcssedForReservation will delete the reservation
func (m *testDBRepo) UpdateProcssedForReservation(id, processed int) error {
	return nil
}

//AllRooms will delete the reservation
func (m *testDBRepo) AllRooms() ([]modals.Room, error) {
	rooms := []modals.Room {
		{ID:1,RoomName:"room1",CreatedAt:time.Now(),UpdatedAt:time.Now()},
		// {ID:2,RoomName:"room2",CreatedAt:time.Now(),UpdatedAt:time.Now()},
	}
	return rooms, nil
}

//GetRestrictionsForRoomByDate will give room restriction by id na ddates
func (m *testDBRepo) GetRestrictionsForRoomByDate(id int, start, end time.Time) ([]modals.RoomRestriction, error) {
	restrictions := []modals.RoomRestriction{
		{ID:1,RoomID:1,ReservationID:1,RestrictionID:1},
		{ID:2,RoomID:2,ReservationID:0,RestrictionID:2},
	}
	return restrictions, nil
}

//InsertBlockForRoom insert the block for room
func (m *testDBRepo) InsertBlockForRoom(id int, start time.Time) error {
	return nil
}

//DeleteBlockByID delete the block by id
func (m *testDBRepo) DeleteBlockByID(id int) error {
	return nil
}
