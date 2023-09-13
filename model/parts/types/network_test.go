package types

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGenerateAndLoad(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	bits := 16
	bin := 8
	size := 10000
	network := &Network{Bits: bits, Bin: bin}
	nodes := network.Generate(size, true)

	filename := fmt.Sprintf("nodes_data_%d_%d.txt", bin, size)
	network.Dump(filename)

	network2 := Network{}
	bits2, bin2, nodes2 := network2.Load(filename)

	//Check if bits2, bin2, nodes2 are the same as bits, bin, nodes
	if bits2 != bits {
		t.Error("Bits are different")
	}
	if bin2 != bin {
		t.Error("Bin are different")
	}
	if len(nodes2) != len(nodes) {
		t.Error("NodesMap are different")
	}
}

//func TestChoice(t *testing.T) {
//	// List of nodes
//	nodes := []NodeId{2, 3, 4, 5, 6, 7, 8, 9, 10}
//	// Originators
//	k := 2
//	c := Choice(nodes, k)
//	assert.Equal(t, len(c), k)
//}
