package main

type node struct {
	pos
	availables []node
}

type path struct {
	nodeA node
	nodeB node
	cost  float32
}

func createPath(pointA, pointB node) path {
	return path{nodeA: pointA, nodeB: pointB, cost: Distance(pointA.pos, pointB.pos)}
}
