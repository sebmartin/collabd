package models

import "gorm.io/gorm"

type Participant struct {
	*gorm.Model
	Name      string
	SessionID int
	Session   Session
}
