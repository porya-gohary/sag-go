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
	"strings"
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
	-e FILE, --precedence FILE   jobset's precedence file
	-n, --naive                  use the naive exploration method [default: false]
	-p, --por                    use the partial-order reduction [default: false]
	-d, --dense-time             use dense time model [default: false]
	-c, --csv                    store the best- and worst-case response times to csv file [default: false]
	-r N, --verbose N            print log messages (0-5) [default: 0]
	-v, --version                show version and exit
	-h, --help                   show this message
`

	arguments, _ := docopt.ParseArgs(argUsage, nil, "0.8.2")

	//Parsing the command-line arguments
	beNaive, _ := arguments.Bool("--naive")
	por, _ := arguments.Bool("--por")
	inputFile, _ := arguments.String("--jobset")
	precedenceFile, _ := arguments.String("--precedence")
	verboseLevel, _ := arguments.Int("--verbose")
	denseTime, _ := arguments.Bool("--dense-time")
	wantCsv, _ := arguments.Bool("--csv")

	start := time.Now()

	commonLogger := verbose.New("Common")
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

	commonLogger.AddHandler("1", sh)
	var workload comm.JobSet
	var csvOutputFile string

	//read job set
	fileExtension := filepath.Ext(inputFile)
	if fileExtension == ".csv" {
		workload = comm.ReadJobSet(inputFile, commonLogger)
	} else if fileExtension == ".yaml" {
		workload = comm.ReadJobSetYAML(inputFile, commonLogger)
	} else {
		commonLogger.Critical("Error: Invalid file extension")
	}

	//read precedence file
	precedenceFileExtension := filepath.Ext(precedenceFile)
	if precedenceFileExtension == ".csv" {
		comm.ReadPrecedence(precedenceFile, &workload, commonLogger)
	} else {
		commonLogger.Critical("Error: Invalid file extension")
	}

	if denseTime {
		comm.WantDenseTimeModel()
	}

	if wantCsv {
		csvOutputFile = strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + ".rta.csv"
	}
	dotOutputFile := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))

	if beNaive {
		if por {
			analysisLogger := verbose.New("NP::Uni::Naive::POR")
			analysisLogger.AddHandler("1", sh)
			uni_non_preemptive_por.ExploreNaively(workload, 10, true, 10, analysisLogger)
			uni_non_preemptive_por.PrintResponseTimes()
			uni_non_preemptive_por.MakeDotFile(dotOutputFile)
			if wantCsv {
				uni_non_preemptive_por.WriteResponseTimes(csvOutputFile)
			}
		} else {
			analysisLogger := verbose.New("NP::Uni::Naive")
			analysisLogger.AddHandler("1", sh)
			uni_non_preemptive.ExploreNaively(workload, 10, true, 10, analysisLogger)
			uni_non_preemptive.PrintResponseTimes()
			uni_non_preemptive.MakeDotFile(dotOutputFile)
			if wantCsv {
				uni_non_preemptive.WriteResponseTimes(csvOutputFile)
			}
		}
	} else {
		if por {
			analysisLogger := verbose.New("NP::Uni::POR")
			analysisLogger.AddHandler("1", sh)
			uni_non_preemptive_por.Explore(workload, 10, true, 10, analysisLogger)
			uni_non_preemptive_por.PrintResponseTimes()
			uni_non_preemptive_por.MakeDotFile(dotOutputFile)
			if wantCsv {
				uni_non_preemptive_por.WriteResponseTimes(csvOutputFile)
			}
		} else {
			analysisLogger := verbose.New("NP::Uni")
			analysisLogger.AddHandler("1", sh)
			uni_non_preemptive.Explore(workload, 10, true, 10, analysisLogger)
			uni_non_preemptive.PrintResponseTimes()
			uni_non_preemptive.MakeDotFile(dotOutputFile)
			if wantCsv {
				uni_non_preemptive.WriteResponseTimes(csvOutputFile)
			}
		}
	}

	fmt.Println("Exploration finished")
	fmt.Println("Time elapsed: ", time.Since(start))

}
