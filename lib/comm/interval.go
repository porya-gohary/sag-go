package comm

type Interval struct {
	Start Time
	End   Time
}

func (i Interval) Intersects(j Interval) bool {
	return i.Start <= j.End && j.Start <= i.End
}

func (i Interval) Contains(j Interval) bool {
	return i.Start <= j.End && j.Start <= i.End
}

func (i Interval) Overlaps(j Interval) bool {
	return i.Start <= j.Start && j.Start <= i.End || i.Start <= j.End && j.End <= i.End
}

func (i Interval) From() Time {
	return i.Start
}

func (i Interval) Min() Time {
	return i.Start
}

func (i Interval) Until() Time {
	return i.End
}

func (i Interval) Max() Time {
	return i.End
}

func (i Interval) Length() int {
	return int(i.End - i.Start)
}

func (i Interval) String() string {
	return "I[" + i.Start.String() + "," + i.End.String() + "]"
}

func (i Interval) Widen(other Interval) Interval {
	return Interval{Start: Minimum(i.Start, other.Start), End: Maximum(i.End, other.End)}
}
