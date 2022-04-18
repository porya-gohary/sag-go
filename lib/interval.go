package lib

import (
	
)

type Interval struct {
	A Time
	B Time
}


func (i Interval) Intersects(j Interval) bool {
	return i.A <= j.B && j.A <= i.B
}

func (i Interval) Contains(j Interval) bool {
	return i.A <= j.B && j.A <= i.B
}

func (i Interval) Overlaps(j Interval) bool {
	return i.A <= j.A && j.A <= i.B || i.A <= j.B && j.B <= i.B
}

func (i Interval) From() Time {
	return i.A
}

func (i Interval) Min() Time {
	return i.A
}

func (i Interval) Until() Time {
	return i.B
}

func (i Interval) Max() Time {
	return i.B
}

func (i Interval) Length() int {
	return int(i.B - i.A)
}

func (i Interval) String() string {
	return "I[" + i.A.String() + "," + i.B.String() + "]"
}

