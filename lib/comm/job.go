package comm

import (
	"fmt"
	"sort"
	"sync"
)

type Job struct {
	Name         string
	TaskID       uint
	JobID        uint
	Arrival      Interval
	Cost         Interval
	Deadline     Time
	Priority     Time
	Predecessors []string
}

type JobSet []*Job

type JobQueue struct {
	queue []*Job
	lock  sync.RWMutex
}

func (j Job) String() string {
	return j.Name + "\t" + j.Arrival.String() + "\t" + j.Cost.String() + "\t" + j.Deadline.String() + "\t" + j.Priority.String() + "\t" + fmt.Sprint(j.Predecessors)
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

func (j Job) GetPredecessors() []string {
	return j.Predecessors
}

func (j Job) PriorityExceeds(otherPriority Time) bool {
	return j.Priority < otherPriority
}

func (j Job) ExceedsDeadline(now Time) bool {
	return (j.Deadline < now) && (now-j.Deadline > DeadlineMissTolerance())
}

func (j *Job) AddPredecessor(predecessor string) {
	j.Predecessors = append(j.Predecessors, predecessor)
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
			return S[i].Name < S[j].Name
		}
		return S[i].Arrival.Start < S[j].Arrival.Start
	})
	return S
}

func (S JobSet) SortByLatestArrival() JobSet {
	sort.Slice(S, func(i, j int) bool {
		if S[i].Arrival.End == S[j].Arrival.End {
			return S[i].Name < S[j].Name
		}
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
		return S[i].HigherPriorityThan(*S[j])
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

func (S *JobSet) GetByName(name string) *Job {
	for _, j := range *S {
		if j.Name == name {
			return j
		}
	}
	return nil
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

func (S JobSet) ContainsByNames(names []string) bool {
	if len(names) == 0 {
		return true
	}
	for _, name := range names {
		if !S.ContainsByName(name) {
			return false
		}
	}
	return true
}

func (S JobSet) ContainsByName(name string) bool {
	for _, j := range S {
		if j.Name == name {
			return true
		}
	}
	return false
}

func (S JobSet) Empty() bool {
	return len(S) == 0
}

//SelectJobByReleaseOrder find job with the lowest release time
func (S *JobSet) SelectJobByReleaseOrder() *Job {
	var job Job
	for _, j := range *S {
		if job.GetEarliestArrival() == 0 {
			job = *j
		} else if job.GetEarliestArrival() > j.GetEarliestArrival() {
			job = *j
		}
	}
	return &job
}

// SelectJobByPriority find job with the lowest priority
func (S *JobSet) SelectJobByPriority() *Job {
	var job Job
	for _, j := range *S {
		if job.HigherPriorityThan(*j) {
			job = *j
		}
	}
	return &job
}

func (S *JobSet) TopologicalSort() JobSet {
	var sortedJobs JobSet
	var temp JobQueue

	for _, job := range *S {
		if job.Predecessors == nil {
			sortedJobs = append(sortedJobs, job)
		} else {
			temp.Enqueue(job)
		}
	}

	for !temp.Empty() {
		job, _ := temp.Front()
		temp.Dequeue()
		predJobs := job.GetPredecessors()

		// All predecessors of j have been sorted already
		if sortedJobs.ContainsByNames(predJobs) {
			sortedJobs = append(sortedJobs, job)
		} else {
			temp.Enqueue(job)
		}
	}

	return sortedJobs
}

func (S *JobSet) SetArrivalTimeWithPrecedence() {
	for _, job := range *S {
		maxEarliestArrival := job.GetEarliestArrival()
		maxLatestArrival := job.GetLatestArrival()
		for _, predJob := range job.GetPredecessors() {
			predJob := S.GetByName(predJob)
			if predJob.GetEarliestArrival() > maxEarliestArrival {
				maxEarliestArrival = predJob.GetEarliestArrival()
			}
			if predJob.GetLatestArrival() > maxLatestArrival {
				maxLatestArrival = predJob.GetLatestArrival()
			}
		}
		job.Arrival.Start = maxEarliestArrival
		job.Arrival.End = maxLatestArrival

	}
}

func (S *JobSet) PreprocessJobs() {
	S.TopologicalSort()
	S.SetArrivalTimeWithPrecedence()
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
