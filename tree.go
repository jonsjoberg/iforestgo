package iforestgo

import (
	"math"
	"math/rand"
)

const EulersConstant = 0.5772156649

type Tree[V Value] struct {
	Root        *Node[V]
	HeightLimit int
}


type Node[V Value] struct {
	Size         int
	SplitPoint   V
	SplitAttrIdx int
	Height       int
	NodeLeft     *Node[V]
	NodeRight    *Node[V]
	External     bool
}

func NewTree[V Value](X *[][]V, sampleIdxs []int, r *rand.Rand) *Tree[V] {
	l := int(math.Ceil(math.Log2(float64(len(sampleIdxs)))))

	return &Tree[V]{
		Root:        nextNode(X, sampleIdxs, 0, l, r),
		HeightLimit: l,
	}
}

func nextNode[V Value](X *[][]V, idxs []int, treeHeight int, heightLimit int, r *rand.Rand) *Node[V] {

	var node Node[V]

	if treeHeight >= heightLimit || len(idxs) <= 1 {
		node = Node[V]{
			Size:     len(idxs),
			Height:   treeHeight,
			External: true,
		}
	} else {

		nAttributes := len((*X)[0])
		// random attribute
		q := r.Intn(nAttributes)

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
	
	node := t.Root
	for {
		if node.External && node.Size == 1 {
			return float64(node.Height)
		} else if node.External {
			return float64(node.Height) + avgPathLength(node.Size)
		}

		if x[node.SplitAttrIdx] < node.SplitPoint {
			node = node.NodeLeft
		} else {
			node = node.NodeRight
		}
	}
}