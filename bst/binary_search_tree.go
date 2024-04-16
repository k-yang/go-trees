package bst

import "fmt"

type Node struct {
	Value int
	Level int
	Left  *Node
	Right *Node
}

func (n *Node) Insert(value int) {
	if n == nil {
		panic("Cannot insert into nil node")
	}

	if n.Value == value {
		return
	}

	if value < n.Value {
		if n.Left == nil {
			n.Left = &Node{Value: value, Level: n.Level + 1}
		} else {
			n.Left.Insert(value)
		}
	} else {
		if n.Right == nil {
			n.Right = &Node{Value: value, Level: n.Level + 1}
		} else {
			n.Right.Insert(value)
		}
	}
}

func (n *Node) Exists(value int) bool {
	if n == nil {
		return false
	}

	if n.Value == value {
		return true
	}

	if value < n.Value {
		return n.Left.Exists(value)
	} else {
		return n.Right.Exists(value)
	}
}

func (n *Node) Traverse() {
	if n == nil {
		return
	}

	n.Left.Traverse()
	fmt.Println(fmt.Sprintf("Value: %d, Level: %d", n.Value, n.Level))
	n.Right.Traverse()
}
