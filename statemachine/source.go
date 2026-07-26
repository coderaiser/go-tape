package statemachine

// TransitionSource loads transition definitions once at startup.
type TransitionSource interface {
	Load() ([]TransitionDef, error)
}

// MemorySource holds transitions in memory.
type MemorySource struct {
	Defs []TransitionDef
}

func (s *MemorySource) Load() ([]TransitionDef, error) {
	return s.Defs, nil
}
