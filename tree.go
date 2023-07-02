package iforestgo

import (
	"fmt"
	"math"
	"math/rand"
)

const EulersConstant = 0.5772156649

type Tree[V Value] struct {
	Root          *Node[V]
	HeightLimit   int
}

func (t *Tree[V]) Print() {
	fmt.Printf("heightLimit: %d\n", t.HeightLimit)
	t.Root.Print()
}

type Node[V Value] struct {
	Size         int
	SplitPoint   V
	SplitAttrIdx int
	NodeLeft     *Node[V]
	NodeRight    *Node[V]
	External     bool
}

func (n *Node[V]) Print() {
	if n.External {
		fmt.Printf("external - size: %d\n", n.Size)
	} else {
		fmt.Printf(
			"splitAttrIdx: %d | splitPoint: %f | size: %d\n",
			n.SplitAttrIdx,
			n.SplitPoint,
			n.Size,
		)
		n.NodeLeft.Print()
		n.NodeRight.Print()
	}
}

func NewTree[V Value](X *[][]V, sampleIdxs []int, r *rand.Rand) *Tree[V] {
	l := int(math.Ceil(math.Log2(float64(len(*X)))))

	return &Tree[V]{
		Root:          nextNode(X, sampleIdxs, 0, l, r),
		HeightLimit:   l,
	}
}

func nextNode[V Value](X *[][]V, idxs []int, treeHeight int, heightLimit int, r *rand.Rand) *Node[V] {

	var node Node[V]

	if treeHeight >= heightLimit || len(idxs) <= 1 {
		node = Node[V]{
			Size:     len(idxs),
			External: true,
		}
	} else {
		// random attribute
		q := r.Intn(len((*X)[0]))

		p := selectSplitPoint(X, idxs, q, r)

		var IdxsL []int
		var IdxsR []int
		for _, i := range idxs {
			v := (*X)[i][q]
			if v < p {
				IdxsL = append(IdxsL, i)
			} else {
				IdxsR = append(IdxsR, i)
			}
		}

		node = Node[V]{
			SplitPoint:   p,
			SplitAttrIdx: q,
			NodeLeft:     nextNode(X, IdxsL, treeHeight+1, heightLimit, r),
			NodeRight:    nextNode(X, IdxsR, treeHeight+1, heightLimit, r),
			External:     false,
		}
	}
	return &node
}

func selectSplitPoint[V Value](X *[][]V, idxs []int, q int, r *rand.Rand) V {

	min := (*X)[idxs[0]][q]
	max := min
	for _, i := range idxs[1:] {
		v := (*X)[i][q]
		if v < min {
			min = v
		} else if v > max {
			max = v
		}
	}

	rF64 := r.Float64()
	return min + V(rF64)*(max-min)
}

func avgPathLength(n int) float64 {
	nf64 := float64(n)
	return 2*(math.Log(nf64-1)+EulersConstant) - ((2 * (nf64 - 1)) / nf64)
}

func PathLength[V Value](x []V, t *Tree[V]) float64 {
	return pathLengthRec[V](x, t.Root, 0)

}

func pathLengthRec[V Value](x []V, n *Node[V], collector int) float64 {
	if n.External {
		if n.Size <= 1 {
			return float64(collector)
		} else {
			return float64(collector) + avgPathLength(n.Size)
		}
	}

	if x[n.SplitAttrIdx] < n.SplitPoint {
		return pathLengthRec[V](x, n.NodeLeft, collector+1)
	} else {
		return pathLengthRec[V](x, n.NodeRight, collector+1)
	}
}

