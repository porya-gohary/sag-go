package comm

import (
	"fmt"
	"sort"
	"sync"
)

type Job struct {
	Name     string
	TaskID   uint
	JobID    uint
	Arrival  Interval
	Cost     Interval
	Deadline Time
	Priority Time
}

type JobSet []*Job

type JobQueue struct {
	queue []*Job
	lock  sync.RWMutex
}

func (j Job) String() string {
	return j.Name + "\t" + j.Arrival.String() + "\t" + j.Cost.String() + "\t" + j.Deadline.String() + "\t" + j.Priority.String()
}

func (j Job) HigherPriorityThan(other Job) bool {

	if j.Priority < other.Priority {
		return true
	}

	if j.Priority == other.Priority {
		// first tie-break by task ID
		if j.TaskID < other.TaskID {
			return true
		} else if j.TaskID == other.TaskID {
			// second, tie-break by job instance
			if j.JobID < other.JobID {
				return true
			}
		}
	}

	return false
}

func (j Job) SameJob(other Job) bool {
	return j.Name == other.Name
}

func (j Job) GetLeastCost() Time {
	return j.Cost.From()
}

func (j Job) GetMaximalCost() Time {
	return j.Cost.Until()
}

func (j Job) GetEarliestArrival() Time {
	return j.Arrival.From()
}

func (j Job) GetLatestArrival() Time {
	return j.Arrival.Until()
}

func (j Job) PriorityExceeds(otherPriority Time) bool {
	return j.Priority < otherPriority
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

func (j JobSet) AbstractString() string {
	var s string
	for _, job := range j {
		s += job.Name + " - "
	}
	return s
}

func (S JobSet) SortByEarliestArrival() JobSet {
	sort.Slice(S, func(i, j int) bool {
		if S[i].Arrival.Start == S[j].Arrival.Start {
			return S[i].Arrival.End < S[j].Arrival.End
		}
		return S[i].Arrival.Start < S[j].Arrival.Start
	})
	return S
}

func (S JobSet) SortByLatestArrival() JobSet {
	sort.Slice(S, func(i, j int) bool {
		return S[i].Arrival.End < S[j].Arrival.End
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

func (S JobSet) SortByWCET() JobSet {
	sort.Slice(S, func(i, j int) bool {
		return S[i].GetMaximalCost() < S[j].GetMaximalCost()
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
		if (*j).Name == job.Name {
			return i
		}
	}
	return -1
}

func (S JobSet) Compare(other JobSet) bool {
	if len(S) != len(other) {
		return false
	}

	sortedJobs := S.SortByEarliestArrival()
	otherSortedJobs := other.SortByEarliestArrival()

	for i, j := range sortedJobs {
		if j.Name != otherSortedJobs[i].Name {
			return false
		}
	}
	return true
}

func (S JobSet) Contains(job Job) bool {
	for _, j := range S {
		if j.Name == job.Name {
			return true
		}
	}
	return false
}

////////////////////////////////
// Functions for job queue

func (c *JobQueue) Enqueue(j *Job) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.queue = append(c.queue, j)
}

func (c *JobQueue) Dequeue() error {
	if len(c.queue) > 0 {
		c.lock.Lock()
		defer c.lock.Unlock()
		c.queue = c.queue[1:]
		return nil
	}
	return fmt.Errorf("Pop Error: Queue is empty")
}

func (c *JobQueue) Front() (*Job, error) {
	if len(c.queue) > 0 {
		c.lock.Lock()
		defer c.lock.Unlock()
		return c.queue[0], nil
	}
	return nil, fmt.Errorf("Peep Error: Queue is empty")
}

func (c *JobQueue) Size() int {
	return len(c.queue)
}

func (c *JobQueue) Empty() bool {
	return len(c.queue) == 0
}
