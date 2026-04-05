package models

import "github.com/google/uuid"

type App struct {
	Id     uuid.UUID
	Name   string
	Secret string
}
