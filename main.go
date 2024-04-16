package main

import "github.com/k-yang/go-trees/avl"

func main() {
	root := &avl.Node{Value: 1}

	root = root.Insert(2)
	root = root.Insert(3)
	root = root.Insert(4)
	root = root.Insert(5)
	root = root.Insert(6)
	root = root.Insert(7)

	root.InOrderTraversal()

	root.BFS()
}
