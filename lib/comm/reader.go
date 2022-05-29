package comm

import (
	"encoding/csv"
	"fmt"
	"github.com/lfkeitel/verbose"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
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

func ReadJobSetYAML(filename string, v *verbose.Logger) JobSet {

	type yamlFile struct {
		Jobset []struct {
			TaskID     uint `yaml:"Task ID"`
			JobID      uint `yaml:"Job ID"`
			ArrivalMin Time `yaml:"Arrival min"`
			ArrivalMax Time `yaml:"Arrival max"`
			CostMin    Time `yaml:"Cost min"`
			CostMax    Time `yaml:"Cost max"`
			Deadline   Time `yaml:"Deadline"`
			Priority   Time `yaml:"Priority"`
		} `yaml:"jobset"`
	}

	var jobs JobSet
	jobSetInYaml := yamlFile{}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	v.Debug("Successfully Opened YAML file")

	if err := yaml.Unmarshal([]byte(file), &jobSetInYaml); err != nil {
		v.Panic(err)
		panic(err)
	}

	for _, job := range jobSetInYaml.Jobset {
		jobInstance := &Job{
			Name:     "J" + fmt.Sprint(job.TaskID) + "," + fmt.Sprint(job.JobID),
			TaskID:   job.TaskID,
			JobID:    job.JobID,
			Arrival:  Interval{Start: job.ArrivalMin, End: job.ArrivalMax},
			Cost:     Interval{Start: job.CostMin, End: job.CostMax},
			Deadline: job.Deadline,
			Priority: job.Priority,
		}
		jobs = append(jobs, jobInstance)
	}

	return jobs
}
