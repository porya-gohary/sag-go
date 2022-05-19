package main

import (
	"fmt"
	"time"
	"go-test/lib"
	"github.com/lfkeitel/verbose"
	
)


func main() {

	start := time.Now()

	//read job set
	workload:=lib.ReadJobSet("./example/fig1c.csv")
	
	// fmt.Println(lib.Infinity())
	// fmt.Println(lib.Epsilon().String())
	


	fmt.Println(time.Since(start))
	logger := verbose.New("app")
	sh := verbose.NewStdoutHandler(true)
	logger.AddHandler("123",sh)



	lib.ExploreNaively(workload, 10, true, 10,logger)

}


