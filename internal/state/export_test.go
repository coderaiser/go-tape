package state

import "github.com/coderaiser/go-tape/internal/statemachine"

var NewFromSource = newFromSource
var ParseTestState = parseTestState
var ParseTestEvent = parseTestEvent

func (s *Store) Adapter() statemachine.Adapter[TestState]     { return s.adapter }
func (s *Store) SetAdapter(a statemachine.Adapter[TestState]) { s.adapter = a }
