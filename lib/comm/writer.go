package comm

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func WriteResponseTimes(filename string, rta map[string]Interval, workload JobSet) {
	csvFile, err := os.Create(filename)
	defer csvFile.Close()
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	w := csv.NewWriter(csvFile)
	defer w.Flush()

	//	write header
	row := []string{"Task ID", "Job ID", "BCCT", "WCCT", "BCRT", "WCRT"}
	if err := w.Write(row); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	//	write data
	for _, j := range workload {
		row := []string{
			fmt.Sprint(j.TaskID),
			fmt.Sprint(j.JobID),
			fmt.Sprint(rta[j.Name].Start.String()),
			fmt.Sprint(rta[j.Name].End.String()),
			fmt.Sprint((rta[j.Name].Start - j.Arrival.Start).String()),
			fmt.Sprint((rta[j.Name].End - j.Arrival.Start).String()),
		}
		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
}
