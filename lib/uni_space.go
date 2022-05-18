package lib

import (
	"fmt"
	"math"
	"time"
)

type resposeTimes map[string]Interval

var beNaive bool = false
var dag *DAG
var states *StateStorage

var statesIndex uint = 0
var currentJobCount int = 0

var startTime time.Time
var elapsedTime time.Duration

var jobsByEarliestArrival JobSet
var jobsByLatestArrival JobSet
var jobsByDeadline JobSet
var jobsByPriority JobSet

// response times
var rta resposeTimes

var aborted bool = false
var deadlineMiss bool = false

func ExploreNaively(workload JobSet, timeout uint, earlyExit bool, maxDepth uint) {
	beNaive = true
	startTime = time.Now()
	explore(workload, timeout, earlyExit, maxDepth)
	elapsedTime = time.Since(startTime)
}

func Explore(workload JobSet, timeout uint, earlyExit bool, maxDepth uint) {
	startTime = time.Now()
	explore(workload, timeout, earlyExit, maxDepth)
	elapsedTime = time.Since(startTime)
}

func explore(workload JobSet, timeout uint, earlyExit bool, maxDepth uint) {
	jobsByEarliestArrival = make(JobSet, len(workload))
	jobsByLatestArrival = make(JobSet, len(workload))
	jobsByDeadline = make(JobSet, len(workload))
	jobsByPriority = make(JobSet, len(workload))

	copy(jobsByEarliestArrival, workload)
	copy(jobsByLatestArrival, workload)
	copy(jobsByDeadline, workload)
	copy(jobsByPriority, workload)

	jobsByEarliestArrival.SortByEarliestArrival()
	jobsByLatestArrival.SortByLatestArrival()
	jobsByDeadline.SortByDeadline()
	jobsByPriority.SortByPriority()

	initialize()

	for currentJobCount < len(workload) {
		frontStates := giveFrontStates()
		for _, state := range frontStates {
			fmt.Println(state.String())
			foundJob:= exploreState(state)

			if !foundJob && len(state.ScheduledJobs) != len(workload) {
				// out of options and we didn't schedule all jobs
				deadlineMiss = true

				if earlyExit{
					aborted = true
					break
				}
			}
		}
		fmt.Println("-----------------------")
		if aborted {
			break
		}

		currentJobCount++
	}

	dag.MakeDot("out")
	

}

func exploreState(state *State) bool {
	var foundJob bool = false
	

	ts_min := state.Availibility.From()
	rel_min := state.EarliestPendingRelease
	t_l := math.Max(float64(nextEligibleJobReady(state)), float64(state.Availibility.Until()))

	nextRange := Interval{Start: Time(math.Min(float64(ts_min), float64(rel_min))), End: Time(t_l)}

	// Iterate over all incomplete jobs that are released no later than nextRange.End
	for _,jt := range jobsByEarliestArrival {
		if jt.Arrival.Start < state.EarliestPendingRelease {
			continue
		}

		if isDispatched(state.ScheduledJobs, jt) {
			continue
		}

		if jt.Arrival.Start > nextRange.End {
			break
		}

		if isEligibleSuccessor(state, jt) {
			schedule(state, jt)
			foundJob = true
		}
	}

	return foundJob
}

func initialize() {
	dag = NewDAG()
	states = NewStateStorage()

	// make root state
	s0 := NewState(statesIndex, Interval{Start: 0, End: 0}, JobSet{}, Time(0))
	v1, _ :=dag.AddVertex(s0.GetName())
	s0.ID=v1
	states.AddState(s0)
	
	fmt.Println("Initialize: ", fmt.Sprint(v1))
	statesIndex++


}

func makeState(finishTime Interval, j JobSet, earliestReleasePending Time,
	parentState *State, dispatchedJob Job) {
	s := NewState(statesIndex, finishTime, j, earliestReleasePending)
	newStateID,_:=dag.AddVertex(s.GetName())
	s.ID=newStateID
	states.AddState(s)
	dag.AddEdge(parentState.GetID(), newStateID, dispatchedJob.Name)
	
	statesIndex++
	fmt.Println("Make state: ", s.String())
	fmt.Println("")
}

func giveFrontStates() []*State {
	leaves := dag.GetLeaves()
	var frontStates []*State
	for _,leaf := range leaves {
		state := states.GetState(fmt.Sprint(leaf))
		// fmt.Println(state.String())
		frontStates = append(frontStates, state)
	}
	return frontStates
}

func nextEligibleJobReady(state *State) Time {

	alreadyScheduled := state.ScheduledJobs
	for _, j := range jobsByLatestArrival {

		// not relevant if already scheduled
		if isDispatched(alreadyScheduled, j) {
			continue
		}

		t := math.Max(float64(j.Arrival.Until()), float64(state.Availibility.Until()))

		// TODO: implement later
		// if (iip_eligible(s, j, t)){
		// 	continue
		// }

		if priorityEligible(state, j, Time(t)) {
			return j.Arrival.Until()
		}

	}
	return Infinity()

}

func isDispatched(jobs JobSet, job Job) bool {
	for _, j := range jobs {
		if j.Name == job.Name {
			return true
		}
	}
	return false
}

func priorityEligible(state *State, j Job, at Time) bool {
	return !certainlyReleasedHigherPriorityExists(state, j, at)
}

func certainlyReleasedHigherPriorityExists(state *State, j Job, at Time) bool {
	// ts_min := state.Availibility.From()
	// rel_min := state.EarliestPendingRelease

	for _, jt := range jobsByLatestArrival {
		// Iterare over all incomplete jobs that are certainly released no later than "at"
		if jt.Arrival.Until() > at {
			break
		}

		// skip reference job
		if jt.Name == j.Name {
			continue
		}

		// TODO: implement later
		// ignore jobs that aren't yet ready
		// if (!ready(s, j)){
		// 	continue
		// }

		// check priority
		if jt.Priority > j.Priority {
			return true
		}

	}
	return false

}

func schedule_eligible_successors(s *State, nextRange Interval) bool {
	for _, j := range jobsByEarliestArrival {

		if j.Arrival.Start < s.EarliestPendingRelease {
			continue
		}

		if j.Arrival.Start > nextRange.End {
			break
		}

		if isDispatched(s.ScheduledJobs, j) {
			continue
		}

		if isEligibleSuccessor(s, j) {

			schedule(s, j)

			return true
		}

	}

	return false
}

func isEligibleSuccessor(s *State, j Job) bool {

	if isDispatched(s.ScheduledJobs, j) {
		return false
	}

	// TODO: implement later
	// if !ready(){
	// 	return false
	// }

	t_s := nextEarliestStartTime(s, j)

	if priorityEligible(s, j, t_s) {
		return false
	}

	if !potentiallyNext(s, j) {
		return false
	}

	// TODO: implement later
	// if (!iip_eligible(s, j, t_s)) {
	// 	return false
	// }

	return true

}

func nextEarliestStartTime(s *State, j Job) Time {
	// t_S in paper, see definition 6.
	return Time(math.Max(float64(j.Arrival.Start), float64(s.Availibility.From())))
}

func potentiallyNext(s *State, j Job) bool {
	t_latest := s.Availibility.Until()

	// if t_latest >=  j.earliest_arrival(), then the
	// job is trivially potentially next, so check the other case.

	if t_latest < j.Arrival.Start {
		r := nextCertainJobRelease(s)

		// if something else is certainly released before j and IIP-
		// eligible at the time of certain release, then j can't
		// possibly be next

		if r < j.Arrival.Start {
			return false
		}

	}
	return true
}

func nextCertainJobRelease(s *State) Time {
	alreadyScheduled := s.ScheduledJobs

	for _, j := range jobsByEarliestArrival {

		if j.Arrival.Start < s.Availibility.Start {
			continue
		}

		// not relevant if already scheduled
		if isDispatched(alreadyScheduled, j) {
			continue
		}

		// TODO: implement later
		// If the job is not IIP-eligible when it is certainly
		// released, then there exists a schedule where it doesn't
		// count, so skip it.

		// if (!iip_eligible(s, j, std::max(j.latestArrival(), s.latest_finish_time())))
		//                 continue;

		// It must be priority-eligible when released, too.
		// Relevant only if we have an IIP, otherwise the job is
		// trivially priority-eligible.

		// if (iip.can_block &&
		// 	!priority_eligible(s, j, std::max(j.latestArrival(), s.latest_finish_time())))
		// 	continue;

		// great, this job fits the bill
		return j.Arrival.End

	}
	return Infinity()

}

func schedule(s *State, j Job) {
	finishRange := nextFinishTimes(s, j)

	scheduledJobs := append(s.ScheduledJobs, j)

	if beNaive {
		makeState(finishRange, scheduledJobs, earliestPossibleJobRelease(s, j), s, j)
	}

}

func nextFinishTimes(s *State, j Job) Interval {
	// standard case -- this job is never aborted or skipped
	i := Interval{Start: nextEarliestFinishTime(s, j), End: nextLatestFinishTime(s, j)}

	return i
}

func nextEarliestFinishTime(s *State, j Job) Time {
	earliestStart := nextEarliestStartTime(s, j)

	return Time(earliestStart + j.Cost.Min())
}

func nextLatestFinishTime(s *State, j Job) Time {
	otherCertainStart := nextCertainHigherPriorityJobRelease(s, j)

	// TODO: implement later
	// t_s := nextEarliestStartTime(s, j)
	// iip_latest_start := iip.latest_start(j, t_s, s);

	// t_s'
	// t_L
	ownLatestStart := math.Max(float64(nextEligibleJobReady(s)), float64(s.Availibility.Until()))

	// t_R, t_I
	// TODO: add iip_latest_start later
	lastStartBeforeOther := otherCertainStart - Epsilon()

	latestFinishTime := Time(math.Min(float64(ownLatestStart), float64(lastStartBeforeOther)))

	return latestFinishTime + j.Cost.Max()

}

func nextCertainHigherPriorityJobRelease(s *State, j Job) Time {
	alreadyScheduled := s.ScheduledJobs

	for _, jt := range jobsByLatestArrival {

		if jt.Arrival.End < s.Availibility.Start {
			continue
		}

		// not relevant if already scheduled
		if isDispatched(alreadyScheduled, jt) {
			continue
		}

		if !jt.higherPriorityThan(j) {
			continue
		}

		// great, this job fits the bill
		return j.Arrival.End

	}
	return Infinity()
}

func earliestPossibleJobRelease(s *State, j Job) Time {
	// Iterate over all incomplete jobs in state s
	for _, jt := range jobsByEarliestArrival {

		if jt.Arrival.Start < s.EarliestPendingRelease {
			continue
		}

		// skip if it is already dispatched
		if isDispatched(s.ScheduledJobs, jt) {
			continue
		}

		// skip if it is the one we're ignoring
		if j.SameJob(jt) {
			continue
		}

		// it's incomplete and not ignored => found the earliest
		return jt.Arrival.Start

	}
	return Infinity()
}
