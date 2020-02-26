// Package idgen is responsible for generating IDs for entities
package idgen

import (
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/shanvl/garbage-events-service"
)

const EventIDLen = 14

// CreateEventID generates an ID for an event entity
func CreateEventID() (garbage.EventID, error) {
	id, err := gonanoid.Nanoid(EventIDLen)
	// an error that rand.Read() might return on very rare occasions
	if err != nil {
		return "", err
	}
	return garbage.EventID(id), nil
}
