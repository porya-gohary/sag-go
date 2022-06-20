package comm

import (
	"fmt"
	"math"
)

var denseTimeModel bool = false

type DiscreteTime int32

type DenseTime float32

type Time float32

func WantDenseTimeModel() {
	denseTimeModel = true
}

func (t Time) String() string {
	if denseTimeModel {
		return fmt.Sprintf("%f", DenseTime(t))
	} else {
		return fmt.Sprintf("%d", DiscreteTime(t))
	}
}

func Infinity() Time {
	if denseTimeModel {
		return Time(math.MaxFloat32)
	} else {
		return Time(math.MaxInt32)
	}
}

func Epsilon() Time {
	if denseTimeModel {
		return Time(math.SmallestNonzeroFloat32)
	} else {
		return Time(1)
	}
}

func Maximum(t1, t2 Time) Time {
	if t1 > t2 {
		return t1
	}
	return t2
}

func Minimum(t1, t2 Time) Time {
	if t1 < t2 {
		return t1
	}
	return t2
}

func DeadlineMissTolerance() Time {
	return Time(1e-6)
}
