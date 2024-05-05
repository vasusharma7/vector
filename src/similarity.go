package vector

import (
	"errors"
	"math"
)

type similarity interface {
	threshold() float64
	calc(Vector, Vector) (float64, error)
}

type SquaredEuclidean struct {
}

func (sqe SquaredEuclidean) threshold() float64 {
	return 1000
}

func (sqe SquaredEuclidean) calc(v1 Vector, v2 Vector) (float64, error) {

	if len(v1.Data) < len(v2.Data) {
		return 0, errors.New("Vector dimensions don't match")
	}
	var dist float64 = 0
	for i := range len(v1.Data) {
		dist += math.Pow((v1.Data[i] - v2.Data[i]), 2)
	}
	return dist, nil
}
