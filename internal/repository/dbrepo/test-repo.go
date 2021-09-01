package dbrepo

import (
	"github.com/VishalTanwani/GolangWebApp/internal/modals"
	"time"
	"errors"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

//InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res modals.Reservation) (int, error) {
	if res.RoomID==2{
		return 0,errors.New("some error")
	}
	return 1, nil
}

//InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(res modals.RoomRestriction) error {
	if res.RoomID==3{
		return errors.New("some error")
	}
	return nil
}

//SearchAvailbilityByDates return status of availability of room
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	if roomID==2{
		return false,nil
	}
	if roomID==3{
		return false,errors.New("some error")
	}
	return true, nil
}

//SearchAvailbilityForAllRooms will give all the rooms available in given dates
func (m *testDBRepo) SearchAvailbilityForAllRooms(start, end time.Time) ([]modals.Room, error) {
	var rooms []modals.Room
	sd := start.Format("2006-01-02")
	if "2021-11-01"==sd{
		return []modals.Room{{ID:1,RoomName:"masdbf"}},nil
	}
	return rooms, nil
}

func (m *testDBRepo) GetRoomByID(id int) (modals.Room, error) {
	var room modals.Room
	if id > 2 {
		return room,errors.New("some error")
	}
	return room, nil
}