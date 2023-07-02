package iforestgo

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	go_iforest "github.com/codegaudi/go-iforest"

	"github.com/stretchr/testify/assert"
)

func TestForest(t *testing.T) {
	X := [][]float32{
		{1.0, 2.0, 3.0},
		{1.2, 2.2, 3.2},
		{5.5, 1.0, 3.2},
		{1.0, 2.0, 3.0},
		{1.1, 2.2, 3.2},
		{1.0, 2.1, 3.0},
		{1.2, 2.2, 3.1},
		{1.0, 2.0, 2.9},
		{1.2, 1.9, 3.0},
	}

	t.Run("ser - deser", func(t *testing.T) {
		f := NewForest[float32](X, 2, 3, 2)

		b, err := f.Serialize()
		assert.NoError(t, err)

		f2, err := Deserialize[float32](b)
		assert.NoError(t, err)

		fmt.Println(f)
		fmt.Println(f2)
	})

	t.Run("create forrest", func(t *testing.T) {
		f := NewForest[float32](X, 100, 9, 2)
		for _, v := range X {
			anomalyScore := f.CalculateAnomalyScore(v)
			fmt.Println(anomalyScore)
		}
	})

}

func BenchmarkComp(b *testing.B) {
	nDataPoints := 1_000_000
	nDimensions := 10
	pctOutliers := 0.001
	nOutliers := int(pctOutliers * float64(nDataPoints))

	r := rand.New(rand.NewSource(1))

	X := make([][]float64, nDataPoints)
	for i := 0; i < nDataPoints; i++ {
		X[i] = make([]float64, nDimensions)
	}

	for i := 0; i < nDimensions; i++ {
		mean := r.NormFloat64()
		for j := 0; j < nDataPoints; j++ {
			X[j][i] = r.NormFloat64() + mean
		}
	}

	outlierIdxs := r.Perm(nDataPoints)[:nOutliers]
	for _, o := range outlierIdxs {
		attrChange := r.Intn(nDimensions)
		// fmt.Println(attrChange)
		X[o][attrChange] += r.NormFloat64() * 10
	}

	b.Run("current", func(b *testing.B) {
		startTrain := time.Now()
		f := NewForest[float64](X, 20, 256, 1)
		endTrain := time.Now()
		fmt.Printf("Current - Time to train: %v\n", endTrain.Sub(startTrain))

		startInference := time.Now()
		for _, x := range X {
			f.CalculateAnomalyScore(x)
		}
		endInference := time.Now()
		fmt.Printf("Current - Time for inference: %v\n", endInference.Sub(startInference))

	})

	b.Run("other", func(b *testing.B) {
		startTrain := time.Now()
		f, _ := go_iforest.NewIForest(X, 20, 256)
		endTrain := time.Now()
		fmt.Printf("Other - Time to train: %v\n", endTrain.Sub(startTrain))

		startInference := time.Now()
		for _, x := range X {
			f.CalculateAnomalyScore(x)
		}
		endInference := time.Now()
		fmt.Printf("Other - Time for inference: %v\n", endInference.Sub(startInference))

	})
}

// Current - Time to train: 2.028323017s
// Current - Time for inference: 13.040929013s
// BenchmarkComp/current-12                       1        15069300716 ns/op       809042984 B/op    164377 allocs/op
// Other - Time to train: 3.893984524s
// Other - Time for inference: 8.623387147s
// BenchmarkComp/other-12                         1        12517404948 ns/op       10406701168 B/op        100054724 allocs/op
