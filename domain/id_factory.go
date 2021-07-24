package domain

import "github.com/google/uuid"

type IDFactory interface {
	UUID() string
}

type iDFactory struct {
}

func NewIDFactory() IDFactory {
	return &iDFactory{}
}

func (u *iDFactory) UUID() string {
	return uuid.New().String()
}
