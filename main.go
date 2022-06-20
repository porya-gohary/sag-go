package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/lfkeitel/verbose"
	"go-test/lib/comm"
	uni_non_preemptive "go-test/lib/uni-non-preemptive"
	uni_non_preemptive_por "go-test/lib/uni-non-preemptive-por"
	"os"
	"path/filepath"
	"time"
)

func main() {

	argUsage := `Unofficial implementation of schedule-abstraction graph analysis with GO
	Copyright © 2022 Pourya Gohari

Usage:
	main [-j FILE] [options]
	main -v
	main -h

Options:
	-j FILE, --jobset FILE       jobset file [default: jobset.csv]
	-n, --naive                  use the naive exploration method [default: false]
	-p, --por                    use the partial-order reduction [default: false]
	-d, --dense-time             use dense time model [default: false]
	-r N, --verbose=N            print log messages (0-5) [default: 0]
	-v, --version                show version and exit
	-h, --help                   show this message
`

	arguments, _ := docopt.ParseArgs(argUsage, nil, "0.8.1")

	//Parsing the command-line arguments
	beNaive, _ := arguments.Bool("--naive")
	por, _ := arguments.Bool("--por")
	inputFile, _ := arguments.String("--jobset")
	verboseLevel, _ := arguments.Int("--verbose")
	denseTime, _ := arguments.Bool("--dense-time")

	start := time.Now()

	logger := verbose.New("NP::uni-non-preemptive")
	sh := verbose.NewStdoutHandler(true)
	//Set verbose level
	if verboseLevel == 0 {
		sh.SetMinLevel(verbose.LogLevelCritical)
	} else if verboseLevel == 1 {
		sh.SetMinLevel(verbose.LogLevelError)
	} else if verboseLevel == 2 {
		sh.SetMinLevel(verbose.LogLevelWarning)
	} else if verboseLevel == 3 {
		sh.SetMinLevel(verbose.LogLevelNotice)
	} else if verboseLevel == 4 {
		sh.SetMinLevel(verbose.LogLevelInfo)
	} else if verboseLevel == 5 {
		sh.SetMinLevel(verbose.LogLevelDebug)
	} else {
		fmt.Println("Error: Invalid verbose level")
		os.Exit(1)
	}

	logger.AddHandler("123", sh)
	var workload comm.JobSet

	//read job set
	fileExtension := filepath.Ext(inputFile)
	if fileExtension == ".csv" {
		workload = comm.ReadJobSet(inputFile, logger)
	} else if fileExtension == ".yaml" {
		workload = comm.ReadJobSetYAML(inputFile, logger)
	} else {
		logger.Critical("Error: Invalid file extension")
	}

	if denseTime {
		comm.WantDenseTimeModel()
	}

	if beNaive {
		//uni_non_preemptive.ExploreNaively(workload, 10, true, 10, logger)
		if por {
			uni_non_preemptive_por.ExploreNaively(workload, 10, true, 10, logger)
			uni_non_preemptive_por.PrintResponseTimes()
		} else {
			uni_non_preemptive.ExploreNaively(workload, 10, true, 10, logger)
			uni_non_preemptive.PrintResponseTimes()
		}
	} else {
		if por {
			uni_non_preemptive_por.Explore(workload, 10, true, 10, logger)
			uni_non_preemptive_por.PrintResponseTimes()
		} else {
			uni_non_preemptive.Explore(workload, 10, true, 10, logger)
			uni_non_preemptive.PrintResponseTimes()
		}
	}

	fmt.Println("Exploration finished")
	fmt.Println("Time elapsed: ", time.Since(start))

}
