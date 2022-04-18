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


func (j Job) String() string {
	return j.Name + "\t" + j.Arrival.String()+ "\t" + j.Cost.String()+ "\t" + j.Deadline.String() + "\t" + j.Priority.String()
}