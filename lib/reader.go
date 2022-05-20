package lib

import (
	"encoding/csv"
	"fmt"
	"github.com/lfkeitel/verbose"
	"os"
	"strconv"
)

func ReadJobSet(filename string, v *verbose.Logger) JobSet {
	var jobs JobSet

	csvFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	v.Debug("Successfully Opened CSV file")

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.TrimLeadingSpace = true

	// skip first line
	if _, err := reader.Read(); err != nil {
		v.Panic(err)
		panic(err)
	}

	csvLines, err := reader.ReadAll()
	if err != nil {
		v.Error(err)
	}

	for _, line := range csvLines {

		taskid, _ := strconv.ParseUint(line[0], 10, 32)
		jobid, _ := strconv.ParseUint(line[1], 10, 32)
		arrivalMin, _ := strconv.Atoi(line[2])
		arrivalMax, _ := strconv.Atoi(line[3])
		costMin, _ := strconv.Atoi(line[4])
		costMax, _ := strconv.Atoi(line[5])
		deadline, _ := strconv.Atoi(line[6])
		priority, _ := strconv.Atoi(line[7])
		jobName := "J" + fmt.Sprint(taskid) + "," + fmt.Sprint(jobid)

		jobInstance := &Job{
			Name:     jobName,
			TaskID:   uint(taskid),
			JobID:    uint(jobid),
			Arrival:  Interval{Start: Time(arrivalMin), End: Time(arrivalMax)},
			Cost:     Interval{Start: Time(costMin), End: Time(costMax)},
			Deadline: Time(deadline),
			Priority: Time(priority),
		}
		// fmt.Println(jobInstance.String())
		jobs = append(jobs, jobInstance)
	}

	return jobs
}
