// Package idgen is responsible for generating IDs for entities
package idgen

import (
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/shanvl/garbage/internal/eventssvc"
)

const (
	eventIDLen = 10
	pupilIDLen = 14

	defLen = 10
)

// CreateEventID generates an ID for an event entity
func CreateEventID() (eventssvc.EventID, error) {
	id, err := gen(eventIDLen)
	return eventssvc.EventID(id), err
}

// CreatePupilID generates an ID for a pupil entity
func CreatePupilID() (eventssvc.PupilID, error) {
	id, err := gen(pupilIDLen)
	return eventssvc.PupilID(id), err
}

func gen(l int) (string, error) {
	if l <= 0 {
		l = defLen
	}
	id, err := gonanoid.Nanoid(l)
	// an error that rand.Read() might return on very rare occasions
	if err != nil {
		return "", err
	}
	return id, nil
}
