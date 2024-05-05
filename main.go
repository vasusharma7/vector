package main

import (
	"fmt"
	"vector/src"
)

func main() {

	fmt.Println("Hey !")
    db := vector.DB{Dim: 2, Sim: vector.SquaredEuclidean{}}
	db.Info()
    
    data1 := [][]float64{{2.0, 4.0}, { 4.0, 2.0}, {5.0, 10}, {10,1},{5,21}}
    data0 := [][]float64{{100,100}, {500,200}, {400,2000}, {100,2000}}
   

    val1 := vector.Vector{Data: []float64{2.0,3.0}, Type: 1}
    if err := db.Insert(val1); err != nil {
        fmt.Println(err.Error())
    }

    for _, d := range data1 {
        vec := vector.Vector{Data: d, Type: 0}
        if err := db.Insert(vec); err != nil {
            fmt.Println(err.Error())
        }
    }

    for _, d := range data0 {
        vec := vector.Vector{Data: d, Type: 0}
        if err := db.Insert(vec); err != nil {
            fmt.Println(err.Error())
        }
    }
    
    db.Info()
    
    //vec := vector.Vector{Data: []float64{2,4}}
    if m,err := db.Search(val1); err == nil {
        fmt.Println("=======Search result=======>", m);
    } else {
        fmt.Println(err.Error());
    }
}
