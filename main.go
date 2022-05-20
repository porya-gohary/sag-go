package main

import (
	"github.com/lfkeitel/verbose"
	"go-test/lib"
	"time"
)

func main() {

	start := time.Now()

	logger := verbose.New("NP::uni")
	sh := verbose.NewStdoutHandler(true)
	// sh.SetMinLevel(verbose.LogLevelInfo)
	logger.AddHandler("123", sh)

	//read job set
	workload := lib.ReadJobSet("./example/example3.csv", logger)

	// fmt.Println(lib.Infinity())
	// fmt.Println(lib.Epsilon().String())

	//lib.ExploreNaively(workload, 10, true, 10, logger)
	lib.Explore(workload, 10, true, 10, logger)
	lib.PrintResponseTimes()
	logger.Info("Naive exploration finished")
	logger.Info("Time elapsed: ", time.Since(start))

}
