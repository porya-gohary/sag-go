package lib

import(
	"fmt"
	"math"
)


type Time float32


func (t Time) String() string {
	return fmt.Sprintf("%f",float32(t))
}

func Infinity() Time {
	return Time(math.MaxFloat32)
}

func Epsilon() Time {
	return Time(math.SmallestNonzeroFloat32)
}