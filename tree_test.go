package iforestgo

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTree(t *testing.T) {

	X := [][]float64{
		{1.0, 2.0, 3.0},
		{1.1, 2.1, 3.1},
		{1.5, 2.5, 3.5},
	}

	r := rand.New(rand.NewSource(2))

	t.Run("select split point", func(t *testing.T) {
		idxs := []int{0, 1, 2}

		q := 1

		res := selectSplitPoint[float64](&X, idxs, q, r)

		assert.True(t, res >= 2 && res <= 2.5)

		idxs = []int{0, 1}
		q = 0
		res = selectSplitPoint[float64](&X, idxs, q, r)
		assert.True(t, res >= 1 && res <= 1.1)

	})

	t.Run("new tree", func(t *testing.T) {
		idxs := []int{0, 1, 2}
		
		tree := NewTree[float64](&X, idxs, r)

		fmt.Println(tree.Root)
	})

}

// func TestSelectSplitPoint(t *testing.T) {

// }

// func TestNewTree() {

// }
