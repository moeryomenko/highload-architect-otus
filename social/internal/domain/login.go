package domain

import uuid "github.com/satori/go.uuid"

// Login represents login model.
type Login struct {
	UserID             uuid.UUID
	Nickname, Password string
}
