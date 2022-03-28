package domain

import uuid "github.com/satori/go.uuid"

type Login struct {
	UserID             uuid.UUID
	Nickname, Password string
}
