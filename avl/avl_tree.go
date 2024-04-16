package avl

import "fmt"

type Node struct {
	Value  int
	Height int
	Left   *Node
	Right  *Node
}

func (n *Node) GetHeight() int {
	if n == nil {
		return -1
	}

	return n.Height
}

func (n *Node) Balance() int {
	if n == nil {
		return 0
	}

	return n.Right.GetHeight() - n.Left.GetHeight()
}

func (n *Node) Insert(value int) *Node {
	fmt.Println(fmt.Sprintf("Inserting %d", value))
	if n == nil {
		return &Node{Value: value, Height: 0}
	}

	if n.Value == value {
		return n
	}

	if value < n.Value {
		n.Left = n.Left.Insert(value)
	} else {
		n.Right = n.Right.Insert(value)
	}

	n.Height = 1 + max(n.Left.GetHeight(), n.Right.GetHeight())
	balance := n.Balance()

	if balance > 1 {
		fmt.Println("Rebalancing left")
		if n.Right.Balance() < 0 {
			fmt.Println("Rotating right")
			n.Right = n.Right.RotateRight()
		}

		fmt.Println("Rotating left")
		n = n.RotateLeft()
	}

	if balance < -1 {
		fmt.Println("Rebalancing right")
		if n.Left.Balance() > 0 {
			fmt.Println("Rotating left")
			n.Left = n.Left.RotateLeft()
		}

		fmt.Println("Rotating right")
		n = n.RotateRight()
	}

	n.Height = 1 + max(n.Left.GetHeight(), n.Right.GetHeight())

	return n
}

func (n *Node) RotateRight() *Node {
	if n == nil {
		return nil
	}

	newRoot := n.Left
	n.Left = newRoot.Right
	newRoot.Right = n

	n.Height = 1 + max(n.Left.GetHeight(), n.Right.GetHeight())
	newRoot.Height = 1 + max(newRoot.Left.GetHeight(), newRoot.Right.GetHeight())

	return newRoot
}

func (n *Node) RotateLeft() *Node {
	if n == nil {
		return nil
	}

	newRoot := n.Right
	n.Right = newRoot.Left
	newRoot.Left = n

	n.Height = 1 + max(n.Left.GetHeight(), n.Right.GetHeight())
	newRoot.Height = 1 + max(newRoot.Left.GetHeight(), newRoot.Right.GetHeight())

	return newRoot
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

func (n *Node) InOrderTraversal() {
	if n == nil {
		return
	}

	n.Left.InOrderTraversal()
	fmt.Println(fmt.Sprintf("Value: %d, Height: %d", n.Value, n.Height))
	n.Right.InOrderTraversal()
}

func (n *Node) BFS() {
	if n == nil {
		return
	}

	queue := []*Node{n}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		fmt.Println(fmt.Sprintf("Value: %d, Height: %d", current.Value, current.Height))

		if current.Left != nil {
			queue = append(queue, current.Left)
		}

		if current.Right != nil {
			queue = append(queue, current.Right)
		}
	}
}
