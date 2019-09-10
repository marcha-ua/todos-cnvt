package parser

type Node struct {
	id int
	Tag      string
	Text     string
	Atr []Attribute
	children []*Node
}

type Attribute struct{
	Name string
	Value string
}

func FindById(root *Node, id int) *Node {
	queue := make([]*Node, 0)
	queue = append(queue, root)
	for len(queue) > 0 {
		nextUp := queue[0]
		queue = queue[1:]
		if nextUp.id == id {
			return nextUp
		}
		if len(nextUp.children) > 0 {
			for _, child := range nextUp.children {
				queue = append(queue, child)
			}
		}
	}
	return nil
}