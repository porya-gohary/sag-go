package lib

import(
	"strconv"
)


type Time int


func (t Time) String() string {
	return strconv.Itoa(int(t))
}