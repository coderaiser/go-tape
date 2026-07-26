package statemachine

// Adapter stores current state per entity id.
type Adapter[S any] interface {
	Get(id string) (S, error)
	Set(id string, state S) error
}
