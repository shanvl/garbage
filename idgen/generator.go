// Package idgen is responsible for generating IDs for entities
package idgen

import (
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/shanvl/garbage-events-service"
)

type IDGenerator interface {
	// GenerateEventID generates an ID for an event entity
	GenerateEventID() (garbage.EventID, error)
}

type idGen struct{}

// GenerateEventID generates an ID for an event entity
func (i *idGen) GenerateEventID() (garbage.EventID, error) {
	id, err := gonanoid.Nanoid(14)
	// this error is an error that rand.Read() might return on very rare occasions
	if err != nil {
		return "", err
	}
	return garbage.EventID(id), nil
}

// NewIDGenerator returns an instance of IDGenerator
func NewIDGenerator() IDGenerator {
	return &idGen{}
}
