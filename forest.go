package iforestgo

import (
	"bytes"
	"encoding/gob"
	"errors"

	"math"
	"math/rand"
)

type Value interface {
	float32 | float64
}

type Forest[V Value] struct {
	Trees           []*Tree[V]
	SubSamplingSize int
	rand            *rand.Rand
	InputDimesion   int
}

var ErrSubSamplingSizeToolarge = errors.New("the requested sub-sampling size exceeds the total number of samples in the input data")

func NewForest[V Value](X [][]V, nTrees int, subSamplingSize int, seed int64) (*Forest[V], error) {

	if len(X) < subSamplingSize {
		return nil, ErrSubSamplingSizeToolarge
	}

	r := rand.New((rand.NewSource(seed)))
	forest := Forest[V]{
		Trees:           make([]*Tree[V], nTrees),
		SubSamplingSize: subSamplingSize,
		rand:            r,
		InputDimesion:   len(X[0]),
	}

	for i := 0; i < nTrees; i++ {
		sampleIdxs := rand.Perm(len(X))[:subSamplingSize]
		forest.Trees[i] = NewTree(&X, sampleIdxs, forest.rand)
	}

	return &forest, nil
}

func (f *Forest[V]) CalculateAnomalyScore(x []V) float64 {
	var sumPathLength float64

	for _, t := range f.Trees {
		sumPathLength += PathLength[V](x, t)
	}

	avgPath := sumPathLength / float64(len(f.Trees))
	return math.Pow(2, -avgPath/avgPathLength(int(f.SubSamplingSize)))
}

func (f *Forest[V]) Serialize() (*bytes.Buffer, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(f)
	return &buf, err
}

func Deserialize[V Value](b *bytes.Buffer) (Forest[V], error) {
	dec := gob.NewDecoder(b)
	var f Forest[V]
	err := dec.Decode(&f)
	return f, err
}
