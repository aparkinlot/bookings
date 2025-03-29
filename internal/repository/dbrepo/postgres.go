package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/aparkinlot/Bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
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

	//var room models.Room
	for rows.Next() {
		var room models.Room
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

// Returns a User by ID
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at
				from users where id = $1`

	row := m.DB.QueryRowContext(cntx, query, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}
	return u, nil
}

// update a User in the Database
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name = $1, last_name = $2, email = $3,
				access_level = $4, updated_at = $5`

	_, err := m.DB.ExecContext(cntx, query, u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// Authenticate validates a user
// Passwords should never be stored as literals in the database -> use hashing
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(cntx, `select id, password from users where email = $1`, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, hashedPassword, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	}
	if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// Returns a slice of all reservations
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
		rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
		order by r.start_date asc
	`

	rows, err := m.DB.QueryContext(cntx, query)
	if err != nil {
		return reservations, err
	}

	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// Returns a slice of all reservations
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at,
		rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
		where processed = 0
		order by r.start_date asc
	`

	rows, err := m.DB.QueryContext(cntx, query)
	if err != nil {
		return reservations, err
	}

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// Returns one reservation by ID
func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
		rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
		where r.id = $1
	`

	row := m.DB.QueryRowContext(cntx, query, id)
	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)
	if err != nil {
		return res, err
	}
	return res, nil
}

// update a Reservation in the Database
func (m *postgresDBRepo) UpdateReservation(u models.Reservation) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update reservations set first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5
		where id = $6
	`

	_, err := m.DB.ExecContext(cntx, query, u.FirstName, u.LastName, u.Email, u.Phone, time.Now(), u.ID)
	if err != nil {
		return err
	}
	return nil
}

// deletes a Reservation from the Database
func (m *postgresDBRepo) DeleteReservation(id int) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from reservations where id = $1`

	_, err := m.DB.ExecContext(cntx, query, id)
	if err != nil {
		return err
	}
	return nil
}

// updates a Reservation from unprocessed to processed
func (m *postgresDBRepo) UpdateProcessedForReservation(id, processed int) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set processed = $1 where id = $2`

	_, err := m.DB.ExecContext(cntx, query, processed, id)
	if err != nil {
		return err
	}
	return nil
}
