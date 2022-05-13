package lib

import (
	"strconv"
)

type State struct {
	Index                  uint
	Availibility           Interval
	ScheduledJobs          JobSet
	EarliestPendingRelease Time
}

type StateStorage map[string]*State

// functions for state
func NewState(index uint, finishTime Interval, j JobSet, earliestRelease Time) *State {

	return &State{
		Index:                  index,
		Availibility:           finishTime,
		ScheduledJobs:          j,
		EarliestPendingRelease: earliestRelease,
	}
}

func (s State) GetName() string {
	return "S" + strconv.FormatUint(uint64(s.Index), 10)

}

func (s State) ID() string {
	return s.GetName()
}

func (s State) String() string {
	return s.GetName() + "\n" + s.Availibility.String() + "\n{" + s.ScheduledJobs.AbstractString() + "}\n" + s.EarliestPendingRelease.String()
}

// functions for state storage

func NewStateStorage() *StateStorage {
	return &StateStorage{}
}

func (s *StateStorage) AddState(st *State) {
	(*s)[(*st).ID()] = st
}

func (s *StateStorage) GetState(name string) *State {
	return (*s)[name]
}

func (s StateStorage) String() string {
	var str string
	for _, v := range s {
		str += v.String() + "\n---------\n"
	}
	return str

}
