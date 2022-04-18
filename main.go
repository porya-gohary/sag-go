package main

import (
	"fmt"
	"go-test/lib"
)


func main() {

	// initialize a new graph
	d := lib.NewDAG()

	// init three vertices
	
	s1 := lib.State{Name: "s1", Age: 50}
	v1, _ := d.AddVertex(s1)

	s2 := lib.State{Name: "s2", Age: 50}
	v2, _ := d.AddVertex(s2)

	s3 := lib.State{Name: "s3", Age: 50}
	v3, _ := d.AddVertex(s3)

	s4 := lib.State{Name: "s4", Age: 50}
	v4, _ := d.AddVertex(s4)

	// add the above vertices and connect them with two edges
	_ = d.AddEdge(v1, v2)
	_ = d.AddEdge(v1, v3)
	_ = d.AddEdge(v3, v4)

	// describe the graph

	// make dot file
	d.MakeDot("out")

	i:= lib.Interval{A: lib.Time(1), B: lib.Time(2)}
	fmt.Println(i.String())

	//read job set
	lib.ReadJobSet("./example/example.csv")
	
}


