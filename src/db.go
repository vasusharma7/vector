package vector

import (
	"fmt"
	"os"
)

const ddir = "./data"

type file struct {
	path string
	ptr *os.File
}

type cluster struct {
    id     int
	fptr   file
	center Vector
    pts    int
}

type Vector struct {
	Data []float64
    Type  int
}

// for now, store everything in memory
type DB struct {
	Dim      int
	Sim      similarity
	clusters []cluster
}

func (db DB) Info() {
	fmt.Printf("No. of clusters %d\n", len(db.clusters))
    for _,c := range db.clusters {
        fmt.Printf("Cluster: id = %d ; center: %v, pts = %d \n", c.id, c.center, c.pts);
    }
}
