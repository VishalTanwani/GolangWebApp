package dbrepo

import (
	"context"
	"github.com/VishalTanwani/GolangWebApp/internal/modals"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

//InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res modals.Reservation) (int, error) {

	//if this transaction is taking longer then give time then time out
	var newID int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id"

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}
	return newID, nil
}

//InsertRoomRestriction inserts a room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(res modals.RoomRestriction) error {

	//if this transaction is taking longer then give time then time out
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "insert into room_restriction (start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7)"

	_, err := m.DB.ExecContext(ctx, stmt,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.ReservationID,
		res.RestrictionID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}
	return nil
}

//SearchAvailbilityByDates return status of availability of room
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	//if this transaction is taking longer then give time then time out
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "select count(id) from room_restriction rr where room_id = $1 and $2 < end_date and $3 > start_date"

	var numRows int

	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)

	if err != nil {
		return false, err
	}

	return numRows == 0, nil
}

//SearchAvailbilityForAllRooms will give all the rooms available in given dates
func (m *postgresDBRepo) SearchAvailbilityForAllRooms(start, end time.Time) ([]modals.Room, error) {

	//if this transaction is taking longer then give time then time out
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []modals.Room
	query := "select r.id, r.room_name from rooms r where r.id not in (select rr.room_id from room_restriction rr where $1 < rr.end_date and $2 > rr.start_date)"

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room modals.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, err

}

func (m *postgresDBRepo) GetRoomByID(id int) (modals.Room, error) {
	//if this transaction is taking longer then give time then time out
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room modals.Room
	query := "select id, room_name from rooms where id = $1"

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(&room.ID, &room.RoomName)
	if err != nil {
		return room, err
	}
	return room, nil
}
