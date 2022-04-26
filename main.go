package main

import (
	"fmt"
	"time"
	"go-test/lib"
)


func main() {

	start := time.Now()

	//read job set
	workload:=lib.ReadJobSet("./example/example.csv")
	
	jobsByArrival:= make(lib.JobSet,len(workload))
	jobsByDeadline:= make(lib.JobSet,len(workload))
	jobsByPriority:= make(lib.JobSet,len(workload))


	copy(jobsByArrival, workload)
	copy(jobsByDeadline, workload)
	copy(jobsByPriority, workload)


	jobsByArrival.SortByArrival()
	jobsByDeadline.SortByDeadline()
	jobsByPriority.SortByPriority()


	// initialize a new graph
	d := lib.NewDAG()

	states := lib.NewStateStorage()
	states.Initialize()
	
	// init three vertices
	s1 := lib.NewState(1,lib.Interval{Start: 0, End: 100}, lib.JobSet{workload[1]}, lib.Time(0))
	states.AddState(s1)
	v1, _ := d.AddVertex(s1.GetName())

	s2 := lib.NewState(2,lib.Interval{Start: 0, End: 100}, lib.JobSet{workload[1],workload[3]}, lib.Time(0))
	states.AddState(s2)
	v2, _ := d.AddVertex(s2.GetName())

	s3 := lib.NewState(3,lib.Interval{Start: 0, End: 100}, lib.JobSet{workload[1],workload[3],workload[2]}, lib.Time(0))
	states.AddState(s3)
	d.AddVertex(s3.GetName())


	// add the above vertices and connect them with two edges
	_ = d.AddEdge(v1, v2,"v1->v2")

	

	// describe the graph

	// make dot file
	d.MakeDot("out")

	// fmt.Println(states.String())
	fmt.Print(jobsByArrival.String())
	fmt.Print(jobsByDeadline.String())
	fmt.Print(jobsByPriority.String())

	fmt.Println(time.Since(start))
	
}


