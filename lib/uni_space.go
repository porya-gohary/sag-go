package lib

import (
	"fmt"
	"github.com/lfkeitel/verbose"
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
var workload JobSet

// response times
var rta resposeTimes

var aborted bool = false
var deadlineMiss bool = false

var logger *verbose.Logger

func ExploreNaively(w JobSet, timeout uint, earlyExit bool, maxDepth uint, v *verbose.Logger) {
	beNaive = true
	logger = v
	startTime = time.Now()
	rta = make(resposeTimes)
	workload = w
	explore(workload, timeout, earlyExit, maxDepth)
	elapsedTime = time.Since(startTime)
}

func Explore(w JobSet, timeout uint, earlyExit bool, maxDepth uint, v *verbose.Logger) {
	logger = v
	startTime = time.Now()
	rta = make(resposeTimes)
	workload = w
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
		frontStates := getFrontStates()
		for _, s := range frontStates {
			logger.Debug("==========================================")
			logger.Debug("Looking at: ", s.GetName())
			foundJob := exploreState(s)
			if !foundJob && len(s.ScheduledJobs) != len(workload) {
				// out of options and we didn't schedule all jobs
				deadlineMiss = true

				if earlyExit {
					aborted = true
					break
				}
			}
		}
		if aborted {
			logger.Warning("---> Aborted!")
			break
		}

		currentJobCount++
	}

	dag.MakeDot("out")

}

func exploreState(s *State) bool {
	var foundJob bool = false

	ts_min := s.Availability.From()
	rel_min := s.EarliestPendingRelease
	t_l := Maximum(nextEligibleJobReady(s), s.Availability.Until())

	nextRange := Interval{Start: Minimum(ts_min, rel_min), End: t_l}

	logger.Debug("ts_min: ", ts_min)
	logger.Debug("rel_min: ", rel_min)
	logger.Debug("t_l: ", t_l)
	logger.Debug("Next range: ", nextRange.String())

	// Iterate over all incomplete jobs that are released no later than nextRange.End
	for _, jt := range jobsByEarliestArrival {
		if jt.Arrival.Start < s.EarliestPendingRelease {
			continue
		}

		if isDispatched(s.ScheduledJobs, *jt) {
			continue
		}

		if jt.GetEarliestArrival() > nextRange.Until() {
			break
		}

		logger.Debug("+ ", jt.Name)
		if isEligibleSuccessor(s, *jt) {
			logger.Debug("  --> can be next ")
			schedule(s, *jt)
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

	v1, _ := dag.AddVertex(s0.GetName(), s0.GetLabel())
	s0.ID = v1
	states.AddState(s0)

	statesIndex++

}

func makeState(finishTime Interval, jobs JobSet, earliestReleasePending Time,
	parentState *State, dispatchedJob Job) {

	s := NewState(statesIndex, finishTime, jobs, earliestReleasePending)
	newStateID, _ := dag.AddVertex(s.GetName(), s.GetLabel())
	s.ID = newStateID

	states.AddState(s)

	edgeLabel := dispatchedJob.Name + "\\nDL=" + fmt.Sprint(dispatchedJob.Deadline)
	edgeLabel += "\\nES=" + fmt.Sprint(finishTime.Start-dispatchedJob.Cost.Start) + "\\nLS=" + fmt.Sprint(finishTime.End-dispatchedJob.Cost.End)
	edgeLabel += "\\nEF=" + fmt.Sprint(finishTime.Start) + "\\nLF=" + fmt.Sprint(finishTime.End)
	dag.AddEdge(parentState.GetID(), newStateID, edgeLabel)
	statesIndex++

	logger.Debug("Make state: ", s.GetName())
	logger.Debug("Availability: ", s.Availability.String())
	logger.Debug("Earliest pending release: ", s.EarliestPendingRelease)
	logger.Debug("Scheduled jobs: ", s.ScheduledJobs.AbstractString())
	logger.Debug("----------------------------------------")
}

func getFrontStates() []*State {
	leaves := dag.GetLeaves()
	var frontStates []*State
	for _, leaf := range leaves {
		s := states.GetState(fmt.Sprint(leaf))
		// fmt.Println(state.String())
		frontStates = append(frontStates, s)

	}
	return frontStates
}

func nextEligibleJobReady(state *State) Time {

	alreadyScheduled := state.ScheduledJobs
	for _, jt := range jobsByLatestArrival {

		// not relevant if already scheduled
		if isDispatched(alreadyScheduled, *jt) {
			continue
		}

		t := Maximum(jt.GetLatestArrival(), state.Availability.Until())

		// TODO: implement later
		// if (iip_eligible(s, j, t)){
		// 	continue
		// }

		if priorityEligible(state, *jt, t) {
			return jt.GetLatestArrival()
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

func priorityEligible(s *State, j Job, at Time) bool {
	return !certainlyReleasedHigherPriorityExists(s, j, at)
}

func certainlyReleasedHigherPriorityExists(s *State, j Job, at Time) bool {
	// ts_min := state.Availability.From()
	// rel_min := state.EarliestPendingRelease
	for _, jt := range jobsByLatestArrival {
		// Iterare over all incomplete jobs that are certainly released no later than "at"

		logger.Debug("        - considering ", jt.Name)
		if jt.GetEarliestArrival() < s.EarliestPendingRelease {
			//fmt.Println("        - 1")
			continue
		}

		if jt.GetLatestArrival() > at {
			//fmt.Println("        - 2")
			break
		}

		if isDispatched(s.ScheduledJobs, *jt) {
			//fmt.Println("        - 3")
			continue
		}

		// skip reference job
		if jt.SameJob(j) {
			//fmt.Println("        - 4")
			continue
		}

		// TODO: implement later
		// ignore jobs that aren't yet ready
		// if (!ready(s, j)){
		// 	continue
		// }

		// check priority
		if jt.higherPriorityThan(j) {
			logger.Debug("=> Found higher priority job: ", jt.Name)
			return true
		}

	}
	return false

}

func scheduleEligibleSuccessors(s *State, nextRange Interval) bool {
	for _, jt := range jobsByEarliestArrival {

		if jt.GetEarliestArrival() < s.EarliestPendingRelease {
			continue
		}

		if jt.GetEarliestArrival() > nextRange.End {
			continue
		}

		if isDispatched(s.ScheduledJobs, *jt) {
			continue
		}

		if isEligibleSuccessor(s, *jt) {

			schedule(s, *jt)

			return true
		}

	}

	return false
}

func isEligibleSuccessor(s *State, j Job) bool {

	if isDispatched(s.ScheduledJobs, j) {
		logger.Debug("Job ", j.Name, "   --> already complete")
		return false
	}

	// TODO: implement later
	// if !ready(){
	// 	return false
	// }

	t_s := nextEarliestStartTime(s, j)

	if !priorityEligible(s, j, t_s) {
		logger.Debug("Job ", j.Name, "   --> not priority eligible")
		return false
	}

	if !potentiallyNext(s, j) {
		logger.Debug("Job ", j.Name, "   --> not potentially next")
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
	return Maximum(s.Availability.From(), j.GetEarliestArrival())
}

func potentiallyNext(s *State, j Job) bool {
	t_latest := s.Availability.Until()

	// if t_latest >=  j.earliest_arrival(), then the
	// job is trivially potentially next, so check the other case.

	if t_latest < j.Arrival.Min() {
		r := nextCertainJobRelease(s)

		// if something else is certainly released before j and IIP-
		// eligible at the time of certain release, then j can't
		// possibly be next

		if r < j.Arrival.Min() {
			return false
		}

	}
	return true
}

func nextCertainJobRelease(s *State) Time {
	alreadyScheduled := s.ScheduledJobs

	for _, jt := range jobsByLatestArrival {

		if jt.GetLatestArrival() < s.Availability.Min() {
			continue
		}

		// not relevant if already scheduled
		if isDispatched(alreadyScheduled, *jt) {
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
		return jt.Arrival.End

	}
	return Infinity()

}

func schedule(parentState *State, j Job) {
	alreadyScheduled := make(JobSet, len(parentState.ScheduledJobs))
	copy(alreadyScheduled, parentState.ScheduledJobs)
	finishRange := nextFinishTimes(parentState, j)

	alreadyScheduled = append(alreadyScheduled, &j)

	logger.Debug("Dispatch job: ", j.Name)

	if beNaive {
		makeState(finishRange, alreadyScheduled, earliestPossibleJobRelease(parentState, j), parentState, j)
	} else {
		if !tryToMerge(finishRange, alreadyScheduled, earliestPossibleJobRelease(parentState, j), parentState, j) {
			makeState(finishRange, alreadyScheduled, earliestPossibleJobRelease(parentState, j), parentState, j)
		}
	}

	updateFinishTimes(j, finishRange)

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
	ownLatestStart := Maximum(nextEligibleJobReady(s), s.Availability.Until())

	logger.Debug("own latest start: ", ownLatestStart)

	// t_R, t_I
	// TODO: add iip_latest_start later
	lastStartBeforeOther := otherCertainStart - Epsilon()

	logger.Debug("last start before other: ", lastStartBeforeOther)

	latestFinishTime := Minimum(ownLatestStart, lastStartBeforeOther)

	return latestFinishTime + j.Cost.Max()

}

func nextCertainHigherPriorityJobRelease(s *State, j Job) Time {
	alreadyScheduled := s.ScheduledJobs

	for _, jt := range jobsByLatestArrival {

		if jt.Arrival.End < s.Availability.Start {
			continue
		}

		// not relevant if already scheduled
		if isDispatched(alreadyScheduled, *jt) {
			continue
		}

		if !jt.higherPriorityThan(j) {
			continue
		}

		// great, this job fits the bill

		return jt.Arrival.Max()

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
		if isDispatched(s.ScheduledJobs, *jt) {
			continue
		}

		// skip if it is the one we're ignoring
		if j.SameJob(*jt) {
			continue
		}

		// it's incomplete and not ignored => found the earliest
		return jt.Arrival.Min()

	}
	return Infinity()
}

func tryToMerge(finishTime Interval, j JobSet, earliestReleasePending Time,
	parentState *State, dispatchedJob Job) bool {
	newState := NewState(statesIndex, finishTime, j, earliestReleasePending)
	tempStates := states.getStatesWithSameJobs(j)
	edgeLabel := dispatchedJob.Name + "\\nDL=" + fmt.Sprint(dispatchedJob.Deadline)
	edgeLabel += "\\nES=" + fmt.Sprint(finishTime.Start-dispatchedJob.Cost.Start) + "\\nLS=" + fmt.Sprint(finishTime.End-dispatchedJob.Cost.End)
	edgeLabel += "\\nEF=" + fmt.Sprint(finishTime.Start) + "\\nLF=" + fmt.Sprint(finishTime.End)

	for _, s := range tempStates {
		if s.IsMergePossible(newState) {
			s.Merge(newState)
			dag.UpdateVertexLabel(s.GetID(), s.GetLabel())
			dag.AddEdge(parentState.GetID(), s.GetID(), edgeLabel)
			return true

		}

	}
	return false

}

func updateFinishTimes(j Job, finishTime Interval) {
	// update the finish time of the job

	if _, ok := rta[j.Name]; ok {
		rta[j.Name] = rta[j.Name].widen(finishTime)
	} else {
		rta[j.Name] = finishTime
	}
}

func PrintResponseTimes() {
	fmt.Println("Response times:")
	fmt.Println("Name: I[BCCT,WCCT]")

	for _, j := range workload {
		fmt.Println(j.Name, ": ", rta[j.Name].String())
	}
}
