package lib

import (
	"fmt"
	"time"
	"math"
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
			exploreState(state)
		}
	}

	currentJobCount++

}

func exploreState(state *State) {
	var foundJob bool = false

	ts_min := state.Availibility.From()
	rel_min := state.EarliestPendingRelease
	t_l := math.Max(float64(nextEligibleJobReady(state)), float64(state.Availibility.Until()))
	

	

	nextRange := Interval{Start: Time(math.Min(float64(ts_min), float64(rel_min))), End: Time(t_l)}

}

func initialize() {
	dag = NewDAG()
	states = NewStateStorage()

	// make root state
	s0 := NewState(statesIndex, Interval{Start: 0, End: 0}, JobSet{}, Time(0))
	states.AddState(s0)
	dag.AddVertex(s0.GetName())
	statesIndex++

	fmt.Println(states.String())

}

func makeState(finishTime Interval, j JobSet, earliestReleasePending Time,
	parentState State, dispatchedJob Job) {
	s := NewState(statesIndex, finishTime, j, earliestReleasePending)
	states.AddState(s)
	dag.AddVertex(s.GetName())
	dag.AddEdge(parentState.GetName(), s.GetName(), dispatchedJob.Name)
	statesIndex++
}

func giveFrontStates() []*State {
	leaves := dag.GetLeaves()
	var frontStates []*State
	for leaf := range leaves {
		state := states.GetState(leaf)
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

		if (priorityEligible(state, j, Time(t))) {
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
