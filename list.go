package bufferedlist

import (
	"os"
	"sync"
)

type node struct {
	data []byte
	next *node
}

var nodePool = sync.Pool{
	New: func() any {
		nb := node{
			data: make([]byte, os.Getpagesize()),
			next: nil,
		}
		return &nb
	},
}

func newNode() *node {
	return nodePool.Get().(*node)
}

func popNode(n *node) {
	nodePool.Put(n)
}

type BufferedList struct {
	head *node
	tail *node
}

func New() BufferedList {
	return BufferedList{}
}
