package repository

import (
	"time"

	"github.com/aparkinlot/Bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailibilityByDatesAndRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailibilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
}
