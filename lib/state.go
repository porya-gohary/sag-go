package lib

import (
	"strconv"
	"fmt"
)

type State struct {
	Index                  uint
	Availibility           Interval
	ScheduledJobs          JobSet
	EarliestPendingRelease Time
	ID  				   string	
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

func (s State) GetID() string {
	return s.ID
}

func (s State) String() string {
	return s.GetName() + "\n" + s.Availibility.String() + "\n{" + s.ScheduledJobs.AbstractString() + "}\n" + s.EarliestPendingRelease.String()
}

func (s State) GetLabel() string {
	var t string
	if (s.EarliestPendingRelease == Infinity()) {
		t="\"" + s.GetName() + ":" + fmt.Sprintf("I[%.3f,%.3f]",s.Availibility.Start,s.Availibility.End) + "\\nER=" + "Inf" + "\""
	} else{
		t="\"" + s.GetName() + ":" + fmt.Sprintf("I[%.3f,%.3f]",s.Availibility.Start,s.Availibility.End) + "\\nER=" + fmt.Sprintf("%.3f",s.EarliestPendingRelease) + "\""
	}
	
	return t
}

// functions for state storage

func NewStateStorage() *StateStorage {
	return &StateStorage{}
}

func (s *StateStorage) AddState(st *State) {
	(*s)[(*st).GetName()] = st
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
