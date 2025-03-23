package models

import "time"

// Reservation data -> requests/json : this is temporarary until database is set up
// type Reservation struct {
// 	FirstName string
// 	LastName  string
// 	Email     string
// 	Phone     string
// }

// User model -> database
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLebel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Room model -> database
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restriction model -> database
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservaiton model -> database
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
}

// Room Restrictions model -> database
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}
