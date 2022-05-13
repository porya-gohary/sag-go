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
	

	lib.ExploreNaively(workload, 10, false, 10)


	fmt.Println(time.Since(start))
	
}


