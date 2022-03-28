package domain

import uuid "github.com/satori/go.uuid"

// Gender represents user gender.
type Gender string

const (
	Male   = "male"
	Female = "female"
)

// User represents user domain model.
type User struct {
	ID   uuid.UUID
	Info *Profile
}

// Profile represents user personal information.
type Profile struct {
	FirstName string
	LastName  string
	Age       int
	Gender    Gender
	Interests []string
	City      string
}
