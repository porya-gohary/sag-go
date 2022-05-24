package uni_non_preemptive_por

import (
	"fmt"
	"go-test/lib/comm"
)

type State struct {
	Index                  uint
	Availability           comm.Interval
	ScheduledJobs          comm.JobSet
	EarliestPendingRelease comm.Time
	ID                     string
}

type StateStorage map[string]*State

// functions for state
func NewState(index uint, finishTime comm.Interval, j comm.JobSet, earliestRelease comm.Time) *State {

	return &State{
		Index:                  index,
		Availability:           finishTime,
		ScheduledJobs:          j,
		EarliestPendingRelease: earliestRelease,
	}
}

func (s *State) GetName() string {
	return "S" + fmt.Sprint(s.Index)
}

func (s State) GetID() string {
	return s.ID
}

func (s State) String() string {
	return s.GetName() + "\n" + s.Availability.String() + "\n{" + s.ScheduledJobs.AbstractString() + "}\n" + s.EarliestPendingRelease.String()
}

func (s State) GetLabel() string {
	var t string
	if s.EarliestPendingRelease == comm.Infinity() {
		t = "\"" + s.GetName() + ":" + fmt.Sprintf("I[%.3f,%.3f]", s.Availability.Start, s.Availability.End) + "\\nER=" + "Inf" + "\""
	} else {
		t = "\"" + s.GetName() + ":" + fmt.Sprintf("I[%.3f,%.3f]", s.Availability.Start, s.Availability.End) + "\\nER=" + fmt.Sprintf("%.3f", s.EarliestPendingRelease) + "\""
	}

	return t
}

func (s State) IsMergePossible(other *State) bool {
	// cannot merge without loss of accuracy if the
	// intervals do not overlap
	if !s.Availability.Intersects(other.Availability) {
		return false
	}

	return true
}

func (s *State) Merge(other *State) {
	(*s).Availability = s.Availability.Widen(other.Availability)

}

// functions for state storage

func NewStateStorage() *StateStorage {
	//return &StateStorage{}
	t := make(StateStorage)
	return &t
}

func (s *StateStorage) AddState(st *State) {
	if _, exists := (*s)[(*st).GetName()]; exists {
		logger.Fatal("Duplicate State!")
	}
	(*s)[(*st).GetName()] = st
}

func (s *StateStorage) GetState(name string) *State {
	return (*s)[name]
}

func (s *StateStorage) String() string {
	var str string
	for _, v := range *s {
		str += v.String() + "\n---------\n"
	}
	return str

}

func (s StateStorage) getStatesWithSameJobs(jobs comm.JobSet) []*State {
	var partialStates []*State
	for _, state := range s {
		if state.ScheduledJobs.Compare(jobs) {
			partialStates = append(partialStates, state)
		}
	}
	return partialStates
}
