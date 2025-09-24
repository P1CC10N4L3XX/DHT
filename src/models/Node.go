package models

type Node struct {
	ID        int
	Parent    *Node
	Childrens **Node
	nephews   **Node
	Resources map[int]string
}
