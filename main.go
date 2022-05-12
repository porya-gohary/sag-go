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

	fmt.Print(jobsByArrival.String())
	fmt.Print(jobsByDeadline.String())
	fmt.Print(jobsByPriority.String())

	lib.ExploreNaively(jobsByArrival, 10, false, 10)


	fmt.Println(time.Since(start))
	
}


