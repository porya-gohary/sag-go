package lib

import(
	"strconv"
	"math"
)


type Time int


func (t Time) String() string {
	return strconv.Itoa(int(t))
}

func Infinity() Time {
	return Time(math.MaxInt)
}

func Epsilon() Time {
	return Time(math.MinInt)
}