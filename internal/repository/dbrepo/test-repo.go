package dbrepo

import (
	"errors"
	"log"
	"time"

	"github.com/aparkinlot/Bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// Inserts a reservation into a database
// If an error happens on the user's end, data will not persist involuntarily -> lifetime of 3 seconds to persist a reservation
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2, then fail; otherwise, pass
	if res.RoomID == 2 {
		return 0, errors.New("some error)")
	}
	return 1, nil
}

// Inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	if r.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

// returns true if availability exists for roomID and false otherwise
func (m *testDBRepo) SearchAvailibilityByDatesAndRoomID(start, end time.Time, roomID int) (bool, error) {

	// set up a test time
	layout := "2006-01-02"
	str := "2049-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	// this is our test to fail the query -- specify 2060-01-01 as start
	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return false, errors.New("some error")
	}

	// if the start date is after 2049-12-31, then return false,
	// indicating no availability;
	if start.After(t) {
		return false, nil
	}

	// otherwise, we have availability
	return true, nil
}

func (m *testDBRepo) SearchAvailibilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	// if the start date is after 2049-12-31, then return empty slice,
	// indicating no rooms are available;
	layout := "2006-01-02"
	str := "2049-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return rooms, errors.New("some error")
	}

	if start.After(t) {
		return rooms, nil
	}

	// otherwise, put an entry into the slice, indicating that some room is
	// available for search dates
	room := models.Room{
		ID: 1,
	}
	rooms = append(rooms, room)

	return rooms, nil
}

func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("some error")
	}

	return room, nil
}
