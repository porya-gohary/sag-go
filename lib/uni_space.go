package lib

import (
	"fmt"
	"time"
)

var beNaive bool = false
var dag *DAG
var states *StateStorage

var statesIndex uint = 0

var startTime time.Time
var elapsedTime time.Duration

func ExploreNaively(workload JobSet, timeout uint, earlyExit bool, maxDepth uint) {
	beNaive = true
	startTime = time.Now()
	explore(workload, timeout, earlyExit, maxDepth)
	elapsedTime = time.Since(startTime)
}

func explore(workload JobSet, timeout uint, earlyExit bool, maxDepth uint) {
	initialize()
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
