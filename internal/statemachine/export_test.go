package statemachine

// EmptyTransitionsMachine builds a machine whose transitions map holds a single
// state with no events out, exercising Validate's empty-transitions branch.
func EmptyTransitionsMachine[S comparable]() *Machine[S, S] {
	var zero S
	return &Machine[S, S]{
		transitions: map[S]map[S]S{zero: {}},
	}
}
