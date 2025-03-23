package dbrepo

import (
	"context"
	"time"

	"github.com/aparkinlot/Bookings/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// Inserts a reservation into a database
// If an error happens on the user's end, data will not persist involuntarily -> lifetime of 3 seconds to persist a reservation
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date,
			end_date, room_id, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(cntx, stmt,
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

// Inserts a room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id,
			created_at, updated_at, restriction_id)
			values
			($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(cntx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)
	if err != nil {
		return err
	}
	return nil
}

// returns true if availability exists for roomID and false otherwise
func (m *postgresDBRepo) SearchAvailibilityByDatesAndRoomID(start, end time.Time, roomID int) (bool, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		select
			count(id)
		from
			room_restrictions
		where
			room_id = $1
			and $2 < end_date and $3 > start_date`

	row := m.DB.QueryRowContext(cntx, query, roomID, start, end)

	var numRows int
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

func (m *postgresDBRepo) SearchAvailibilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	query := `
		select
			r.id, r.room_name
		from
			rooms r
		where r.id not in (select room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date)
		`

	rows, err := m.DB.QueryContext(cntx, query, start, end)
	if err != nil {
		return rooms, err
	}

	var room models.Room
	for rows.Next() {
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `
		select id, room_name, created_at, updated_at from rooms where id = $1
	`

	row := m.DB.QueryRowContext(cntx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}
	return room, nil
}
