package lib

type Job struct {
	Name string
	TaskID int
	JobID int
	Arrival Interval
	Cost Interval
	Deadline Time
	Priority Time
}

type JobSet []Job


func (j Job) String() string {
	return j.Name + "\t" + j.Arrival.String()+ "\t" + j.Cost.String()+ "\t" + j.Deadline.String() + "\t" + j.Priority.String()
}

func (j JobSet) String() string {
	var s string
	for _, job := range j {
		s += job.String() + "\n"
	}
	return s
}

func (j JobSet) AbstractString() string{
	var s string
	for _, job := range j {
		s += job.Name + ","
	}
	return s
}