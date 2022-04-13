package main

import (
	"fmt"
	"log"
    "os"
	"github.com/heimdalr/dag"
)

type state struct {
    name string
    age  int
}

func main() {

	// initialize a new graph
	d := dag.NewDAG()

	// init three vertices
	s1 := state {name: "s1", age: 50}
	v1, _ := d.AddVertex(s1)
	
	s2 := state {name: "s2", age: 50}
	v2, _ := d.AddVertex(s2)

	s3 := state {name: "s3", age: 50}
	v3, _ := d.AddVertex(s3)

	s4 := state {name: "s4", age: 50}
	v4, _ := d.AddVertex(s4)

	// add the above vertices and connect them with two edges
	_ = d.AddEdge(v1, v2)
	_ = d.AddEdge(v1, v3)
	_ = d.AddEdge(v3, v4)
	

	// describe the graph

	// make dot file
	f, err := os.Create("out.dot")

	if err != nil {
        log.Fatal(err)
    }

	defer f.Close()

    _, err2 := f.WriteString(makeDot(d))

    if err2 != nil {
        log.Fatal(err2)
    }

    fmt.Println("Done!")
}

func makeDot(d *dag.DAG) string{
	var dotOut string
	dotOut += "digraph {\n"
	dotOut += "\tgraph [fontname=Ubuntu];\n"
	dotOut += "\tnode [fontname=Ubuntu];\n"
	dotOut += "\tedge [fontname=Ubuntu];\n"
	for _, v := range d.GetVertices() {
		dotOut += fmt.Sprintf("\t%v[label=%v];\n",v.(state).name ,v.(state).name)
	}

	for key, v := range d.GetVertices() {
		x,_ := d.GetChildren(key)

		for vertex := range x {
			child,_ :=d.GetVertex(vertex)
			dotOut += fmt.Sprintf("\t%v -> %v;\n",v.(state).name ,child.(state).name)
		}
	}
	dotOut += "}"

	return dotOut
	// fmt.Print(d.String())

}