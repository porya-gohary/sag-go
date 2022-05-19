package main

import (
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
	


	
	logger := verbose.New("NP::uni")
	sh := verbose.NewStdoutHandler(true)
	logger.AddHandler("123",sh)



	lib.ExploreNaively(workload, 10, true, 10,logger)
	logger.Info("Naive exploration finished")
	logger.Info("Time elapsed: ", time.Since(start))

}


