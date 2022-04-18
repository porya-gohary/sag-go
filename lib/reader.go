package lib

import (
    "encoding/csv"
    "fmt"
    "os"
	"strconv"
)


func ReadJobSet(filename string) []Job {

	csvFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened CSV file")

	defer csvFile.Close()


    reader := csv.NewReader(csvFile)

	// skip first line
	if _, err := reader.Read(); err != nil {
		panic(err)
	}

    csvLines, err := reader.ReadAll()
    if err != nil {
        fmt.Println(err)
    }    


    for _, line := range csvLines {

		taskid,_ :=strconv.Atoi(line[0])
		jobid,_ :=strconv.Atoi(line[1])
		arrivalMin,_ :=strconv.Atoi(line[2])
		arrivalMax,_ :=strconv.Atoi(line[3])
		costMin,_ :=strconv.Atoi(line[4])
		costMax,_ :=strconv.Atoi(line[5])
		deadline,_ :=strconv.Atoi(line[6])
		priority,_ :=strconv.Atoi(line[7])
		jobName := "J"+strconv.Itoa(taskid)+","+strconv.Itoa(jobid)

        jobInstance := Job{
			Name: jobName,
			TaskID: taskid,
			JobID: jobid,
			Arrival: Interval{A: Time(arrivalMin), B: Time(arrivalMax)},
			Cost: Interval{A: Time(costMin), B: Time(costMax)},
			Deadline: Time(deadline),
			Priority: Time(priority),
        }
        fmt.Println(jobInstance.String())
    }

	return nil
}