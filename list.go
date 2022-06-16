package bufferedlist

import (
	"os"
	"sync"
)

var maxDataLen = os.Getpagesize()

type node struct {
	data  []byte
	start int
	end   int
	next  *node
}

func (n *node) append(data []byte) int {
	appendingLength := maxDataLen - n.end
	if appendingLength > len(data) {
		appendingLength = len(data)
	}
	for i := n.end; i < n.end+appendingLength; i++ {
		n.data[i] = data[i-n.end]
	}
	n.end += appendingLength
	return appendingLength
}

func (n *node) read(data []byte) (int, bool) {
	readingLen := n.end - n.start
	if readingLen > len(data) {
		readingLen = len(data)
	}
	for i := n.start; i < n.start+readingLen; i++ {
		data[i-n.start] = n.data[i]
	}
	n.start += readingLen
	return readingLen, n.start == n.end
}

var nodePool = sync.Pool{
	New: func() any {
		nb := node{
			data: make([]byte, maxDataLen),
			next: nil,
		}
		return &nb
	},
}

func newNode() *node {
	return nodePool.Get().(*node)
}

func popNode(n *node) {
	n.start = 0
	n.end = 0
	for i := 0; i < len(n.data); i++ {
		n.data[i] = 0
	}
	n.next = nil
	nodePool.Put(n)
}

type BufferedList struct {
	head *node
	tail *node
}

func New() BufferedList {
	return BufferedList{}
}

func (b *BufferedList) Append(data []byte) {
	if b.head == nil {
		b.head = newNode()
		b.tail = b.head
	}
	appendedLen := b.tail.append(data)
	for appendedLen < len(data) {
		b.tail.next = newNode()
		b.tail = b.tail.next
		appendedLen += b.tail.append(data[appendedLen:])
	}
}

func (b *BufferedList) Read(data []byte) int {
	if b.head == nil {
		return 0
	}
	index := 0
	for {
		readLen, isEmpty := b.head.read(data[index:])
		index += readLen
		if isEmpty {
			popNode(b.head)
			b.head = b.head.next
			if b.head == nil {
				b.tail = nil
				break
			}
		}
		if index == len(data) {
			break
		}
	}
	return index
}
