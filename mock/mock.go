package mock

import "github.com/shanvl/garbage-events-service"

// EventingRepository is a mock eventing usecase repository
type EventingRepository struct {
	StoreEventFn      func(e *garbage.Event) (garbage.EventID, error)
	StoreEventInvoked bool
}

// StoreEvent calls StoreEventFn
func (r *EventingRepository) StoreEvent(e *garbage.Event) (garbage.EventID, error) {
	r.StoreEventInvoked = true
	return r.StoreEventFn(e)
}

type IDGenerator struct {
	GenerateEventIDFn      func() garbage.EventID
	GenerateEventIDInvoked bool
}

func (g *IDGenerator) GenerateEventID() garbage.EventID {
	g.GenerateEventIDInvoked = true
	return g.GenerateEventIDFn()
}
