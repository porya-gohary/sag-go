package lib

import (
	"sort"
)

type Job struct {
	Name string
	TaskID int
	JobID int
	Arrival Interval
	Cost Interval
	Deadline Time
	Priority Time
}

type JobSet []Job


func (j Job) String() string {
	return j.Name + "\t" + j.Arrival.String()+ "\t" + j.Cost.String()+ "\t" + j.Deadline.String() + "\t" + j.Priority.String()
}




////////////////////////////////
// Functions for jobset
func (j JobSet) String() string {
	var s string
	for _, job := range j {
		s += job.String() + "\n"
	}
	return s
}

func (j JobSet) AbstractString() string{
	var s string
	for _, job := range j {
		s += job.Name + ","
	}
	return s
}


func (S JobSet) SortByArrival() JobSet {
	sort.Slice(S, func(i, j int) bool {
		return S[i].Arrival.Start < S[j].Arrival.Start
	})
	return S
}


func (S JobSet) SortByDeadline() JobSet {
	sort.Slice(S, func(i, j int) bool {
		return S[i].Deadline < S[j].Deadline
	})
	return S
}

func (S JobSet) SortByPriority() JobSet {
	sort.Slice(S, func(i, j int) bool {
		return S[i].Priority < S[j].Priority
	})
	return S
}

func (S JobSet) RemoveByIndex(index int) JobSet {
	return append(S[:index], S[index+1:]...)
}

func (S JobSet) Remove(job Job) JobSet {
	return S.RemoveByIndex(S.IndexOf(job))
}

func (S JobSet) IndexOf(job Job) int {
	for i, j := range S {
		if j == job {
			return i
		}
	}
	return -1
}