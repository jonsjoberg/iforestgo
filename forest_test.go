package iforestgo

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	go_iforest "github.com/codegaudi/go-iforest"

	"github.com/stretchr/testify/assert"
)

func compTree[V Value](a *Node[V], b *Node[V]) bool {
	if a.External && b.External {
		return a.Size == b.Size
	} else if a.SplitPoint != b.SplitPoint ||
		a.SplitAttrIdx != b.SplitAttrIdx ||
		a.External != b.External {
		return false
	}

	return compTree[V](a.NodeLeft, b.NodeLeft) && compTree[V](a.NodeRight, b.NodeRight)
}

func TestForest(t *testing.T) {
	// idx 2 is the outlier
	X := [][]float64{
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
	t.Run("create forrest", func(t *testing.T) {
		f, err := NewForest[float64](X, 100, 9, 2)
		assert.NoError(t, err)

		maxNormalScore := 0.0
		outlierScore := 0.0

		for i, v := range X {
			anomalyScore := f.CalculateAnomalyScore(v)
			if i == 2 {
				outlierScore = anomalyScore
			} else if anomalyScore > maxNormalScore {
				maxNormalScore = anomalyScore
			}
		}

		assert.Greater(t, outlierScore, maxNormalScore)
	})

	t.Run("sample size to large", func(t *testing.T) {
		f, err := NewForest[float64](X, 100, 100, 2)
		assert.Error(t, ErrSubSamplingSizeToolarge, err)
		assert.Nil(t, f)
	})

	t.Run("ser - deser", func(t *testing.T) {
		f, err := NewForest[float64](X, 2, 3, 2)
		assert.NoError(t, err)

		b, err := f.Serialize()
		assert.NoError(t, err)

		f2, err := Deserialize[float64](b)
		assert.NoError(t, err)

		assert.Equal(t, len(f.Trees), len(f2.Trees))
		assert.Equal(t, f.SubSamplingSize, f2.SubSamplingSize)

		for i, t1 := range f.Trees {
			t2 := f2.Trees[i]
			assert.Equal(t, t1.HeightLimit, t2.HeightLimit)
			assert.True(t, compTree[float64](t1.Root, t2.Root))
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

	b.Run("go_iforest", func(b *testing.B) {
		startTrain := time.Now()
		f, _ := go_iforest.NewIForest(X, 100, 256)
		endTrain := time.Now()
		fmt.Printf("go_iforest - Time to train: %v\n", endTrain.Sub(startTrain))

		startInference := time.Now()
		for _, x := range X {
			f.CalculateAnomalyScore(x)
		}
		endInference := time.Now()
		fmt.Printf("go_iforest - Time for inference: %v\n", endInference.Sub(startInference))

	})

	b.Run("iforest-go", func(b *testing.B) {
		startTrain := time.Now()
		f, _ := NewForest[float64](X, 100, 256, 1)
		endTrain := time.Now()
		fmt.Printf("iforest-go - Time to train: %v\n", endTrain.Sub(startTrain))

		startInference := time.Now()
		for _, x := range X {
			f.CalculateAnomalyScore(x)
		}
		endInference := time.Now()
		fmt.Printf("iforest-go - Time for inference: %v\n", endInference.Sub(startInference))

	})
}
