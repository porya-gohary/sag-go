package uni_non_preemptive_por

import "go-test/lib/comm"

type reductionSet struct {
	jobs                    comm.JobSet
	jobsByEarliestArrival   comm.JobSet
	jobsByLatestArrival     comm.JobSet
	jobsByWCET              comm.JobSet
	latestBusyTime          comm.Time
	latestIdleTime          comm.Time
	latestStartTimes        map[string]comm.Time
	maxPriority             comm.Time
	numInterferingJobsAdded uint
	availability            comm.Interval
}

func createReductionSet(s *State, eligibleSuccessors comm.JobSet) *reductionSet {
	jobsByEarliestArrivalLocal := make(comm.JobSet, len(eligibleSuccessors))
	jobsByLatestArrivalLocal := make(comm.JobSet, len(eligibleSuccessors))
	jobsByWCETLocal := make(comm.JobSet, len(eligibleSuccessors))

	copy(jobsByEarliestArrivalLocal, eligibleSuccessors)
	copy(jobsByLatestArrivalLocal, eligibleSuccessors)
	copy(jobsByWCETLocal, eligibleSuccessors)

	jobsByEarliestArrival.SortByEarliestArrival()
	jobsByLatestArrival.SortByLatestArrival()
	jobsByWCETLocal.SortByWCET()

	rs := &reductionSet{
		jobs:                    eligibleSuccessors,
		jobsByEarliestArrival:   jobsByEarliestArrival,
		jobsByLatestArrival:     jobsByLatestArrival,
		jobsByWCET:              jobsByWCETLocal,
		latestBusyTime:          comm.Time(0),
		latestIdleTime:          comm.Time(0),
		maxPriority:             comm.Time(0),
		numInterferingJobsAdded: 0,
		availability:            s.Availability,
	}

	rs.setLatestBusyTime()
	rs.setLatestIdleTime()
	rs.setLatestStartTimes()
	rs.setMaxPriority()

	return rs

}

func (rs *reductionSet) setLatestBusyTime() {
	t := rs.availability.Max()
	for _, j := range rs.jobsByLatestArrival {
		t = comm.Maximum(t, j.GetLatestArrival()) + j.GetMaximalCost()
	}
	rs.latestBusyTime = t
}

func (rs *reductionSet) setLatestIdleTime() {
	idleTime := comm.Time(-1)

	var idleJob *comm.Job

	for _, i := range rs.jobsByLatestArrival {
		if i.GetLatestArrival() > rs.availability.Min() {
			t := rs.availability.Min()
			for _, j := range rs.jobsByEarliestArrival {
				if j.GetLatestArrival() < i.GetLatestArrival() {
					t = comm.Maximum(t, j.GetEarliestArrival()) + j.GetLeastCost()
				}

				if t >= i.GetLatestArrival() {
					break
				}
			}

			if t < i.GetLatestArrival() {
				if idleJob == nil || i.GetLatestArrival() > idleJob.GetLatestArrival() {
					idleJob = i
				}
			}

		}
	}

	if idleJob == nil {
		rs.latestIdleTime = idleTime
		return
	}

	if idleJob.GetLatestArrival() == rs.jobsByLatestArrival[0].GetLatestArrival() {
		rs.latestIdleTime = idleTime
		return
	} else {
		rs.latestIdleTime = idleJob.GetLatestArrival() - comm.Epsilon()
	}

}

func (rs *reductionSet) setLatestStartTimes() {
	jobsByPrio := rs.preprocessPriorities()
	var startTimes map[string]comm.Time
	for _, j := range rs.jobs {
		startTimes[j.Name] = rs.computeLatestStartTime(j, jobsByPrio)
	}
	rs.latestStartTimes = startTimes
}

// Preprocess priorities for s_i by setting priority of each job to the lowest priority of its predecessors
func (rs *reductionSet) preprocessPriorities() map[string]comm.Time {
	var jobsByPrio map[string]comm.Time
	for _, j := range rs.jobs {
		maxPredPrio := comm.Time(0)
		//	TODO: implement precedence constraints
		p := comm.Maximum(maxPredPrio, j.Priority)
		jobsByPrio[j.Name] = p
	}
	return jobsByPrio
}

func (rs *reductionSet) computeLatestStartTime(j *comm.Job, jobsByPrio map[string]comm.Time) comm.Time {
	s_i := rs.computeSi(j, jobsByPrio)

	return comm.Minimum(s_i, rs.computeSecondLstBound(j))
}

// Upper bound on latest start time (s_i)
func (rs *reductionSet) computeSi(i *comm.Job, jobsByPrio map[string]comm.Time) comm.Time {
	var blockingJob *comm.Job
	var blockingTime comm.Time
	var latestStartTime comm.Time

	for _, j := range rs.jobsByEarliestArrival {
		if i.SameJob(*j) {
			continue
		}

		// use preprocessed prio level
		if i.PriorityExceeds(jobsByPrio[j.Name]) && (blockingJob == nil || blockingJob.GetMaximalCost() < j.GetMaximalCost()) {
			blockingJob = j

		}
	}

	if blockingJob == nil {
		blockingTime = comm.Time(0)
	} else {
		blockingTime = comm.Maximum(0, blockingJob.GetMaximalCost()-comm.Epsilon())
	}
	latestStartTime = comm.Maximum(rs.availability.Max(), i.GetLatestArrival()+blockingTime)

	for _, j := range rs.jobsByEarliestArrival {
		if i.SameJob(*j) {
			continue
		}

		if j.GetEarliestArrival() <= latestStartTime && !i.PriorityExceeds(jobsByPrio[j.Name]) {
			latestStartTime += j.GetMaximalCost()
		} else if j.GetEarliestArrival() > latestStartTime {
			break
		}

	}
	return latestStartTime

}

// Upper bound on latest start time (LFT^bar - sum(C_j^max) - C_i^max)
func (rs *reductionSet) computeSecondLstBound(j *comm.Job) comm.Time {
	descendants := rs.getDescendants(j)
	sum := comm.Time(0)
	for _, d := range descendants {
		sum += d.GetMaximalCost()
	}

	return rs.latestBusyTime - sum - j.GetMaximalCost()
}

// Gets all descendants of J_i in J^M
func (rs *reductionSet) getDescendants(j *comm.Job) []*comm.Job {
	var descendants comm.JobSet
	var remainingJobs comm.JobSet = make(comm.JobSet, len(rs.jobs))
	copy(remainingJobs, rs.jobs)

	var queue comm.JobQueue
	queue.Enqueue(j)
	for !queue.Empty() {
		//jt := queue.Dequeue()
		queue.Dequeue()
		//	TODO: implement precedence constraints
		for _, i := range remainingJobs {
			if descendants.Contains(*i) {
				remainingJobs.Remove(*i)
			}
		}
	}
	return descendants
}

func (rs *reductionSet) setMaxPriority() {
	var maxPriority comm.Time
	for _, j := range rs.jobs {
		if !j.PriorityExceeds(maxPriority) {
			maxPriority = j.Priority
		}
	}
	rs.maxPriority = maxPriority
}