package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type node struct {
	pos pos
	id  int
}

type path struct {
	nodeA node
	nodeB node
	cost  float32
}

func findNodeByID(id int) *node {
	for _, node := range game.currentmap.nodes {
		if node.id == id {
			return &node
		}
	}
	fmt.Printf("Warning: Node with ID %d not found\n", id)
	return nil
}
func createNode(id int, pos pos) node {
	return node{
		id:  id,
		pos: pos,
	}
}

func createPath(nodeA *node, nodeB *node, cost float32) path {
	if nodeA == nil || nodeB == nil {
		fmt.Println("Error: Cannot create path with nil nodes")
		return path{}
	}
	return path{
		nodeA: *nodeA,
		nodeB: *nodeB,
		cost:  cost,
	}
}

func drawPath(s *ebiten.Image, path path) {
	ebitenutil.DrawLine(s,
		float64(offsetsx(path.nodeA.pos.float_x)), float64(offsetsy(path.nodeA.pos.float_y)),
		float64(offsetsx(path.nodeB.pos.float_x)), float64(offsetsy(path.nodeB.pos.float_y)), uidarkred)
}

// Closest point on a line segment from target
func closestPointOnSegment(target, a, b pos) pos {
	// Vector AB
	ab := pos{float_x: b.float_x - a.float_x, float_y: b.float_y - a.float_y}
	// Vector AT
	at := pos{float_x: target.float_x - a.float_x, float_y: target.float_y - a.float_y}

	// Dot product of AB and AT
	dotProduct := ab.float_x*at.float_x + ab.float_y*at.float_y
	// Length squared of AB
	abLenSq := ab.float_x*ab.float_x + ab.float_y*ab.float_y

	// Projection scalar (clamped between 0 and 1)
	t := dotProduct / abLenSq
	if t < 0 {
		t = 0 // Closest to point A
	} else if t > 1 {
		t = 1 // Closest to point B
	}

	// Closest point on the segment
	return pos{
		float_x: a.float_x + t*ab.float_x,
		float_y: a.float_y + t*ab.float_y,
	}
}

func findClosestNode(target pos) node {
	var rn node
	leastDistance := float32(1e9) // Initialize to a very large value

	for _, n := range game.currentmap.nodes {
		if Distance(n.pos, target) < leastDistance {
			leastDistance = Distance(n.pos, target)
			rn = n
		}
	}

	return rn
}

func findClosestPointOnPaths(target pos) (pos, float32) {
	var closestPoint pos
	minDistance := float32(math.MaxFloat32)

	for _, p := range game.currentmap.paths {
		point := closestPointOnSegment(target, p.nodeA.pos, p.nodeB.pos)
		d := Distance(target, point)
		if d < minDistance {
			minDistance = d
			closestPoint = point
		}
	}
	return closestPoint, minDistance
}

func nodesWithinRange(startNode node, maxHops int) map[int]bool {
	visited := make(map[int]bool)
	queue := []struct {
		nodeID int
		hops   int
	}{{nodeID: startNode.id, hops: 0}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.hops > maxHops {
			continue
		}

		if visited[current.nodeID] {
			continue
		}

		visited[current.nodeID] = true

		// Enqueue neighbors
		for _, p := range game.currentmap.paths {
			// Check if the current node is part of the path
			var neighborID int
			if p.nodeA.id == current.nodeID {
				neighborID = p.nodeB.id
			} else if p.nodeB.id == current.nodeID {
				neighborID = p.nodeA.id
			} else {
				continue
			}

			if !visited[neighborID] {
				queue = append(queue, struct {
					nodeID int
					hops   int
				}{nodeID: neighborID, hops: current.hops + 1})
			}
		}
	}

	return visited
}

func randomPointWithinRange(startNode node, maxHops int) pos {
	// Find nodes within range using BFS
	reachableNodes := nodesWithinRange(startNode, maxHops)

	// Collect paths between reachable nodes
	validPaths := []path{}
	for _, p := range game.currentmap.paths {
		if reachableNodes[p.nodeA.id] && reachableNodes[p.nodeB.id] {
			validPaths = append(validPaths, p)
		}
	}

	// If no valid paths exist, return the starting node's position
	if len(validPaths) == 0 {
		return startNode.pos
	}

	// Select a random path
	rand.Seed(time.Now().UnixNano())
	randomPath := validPaths[rand.Intn(len(validPaths))]

	// Generate a random point on the selected path
	t := rand.Float32()
	x := randomPath.nodeA.pos.float_x + t*(randomPath.nodeB.pos.float_x-randomPath.nodeA.pos.float_x)
	y := randomPath.nodeA.pos.float_y + t*(randomPath.nodeB.pos.float_y-randomPath.nodeA.pos.float_y)

	return pos{float_x: x, float_y: y}
}

type priorityQueueItem struct {
	nodeID   int
	cost     float32
	previous *priorityQueueItem
}

func findShortestPathPositions(startID, goalID int) []pos {
	// Priority queue for exploring nodes, ordered by cost
	pq := []priorityQueueItem{
		{nodeID: startID, cost: 0, previous: nil},
	}

	// Map to store the shortest cost to reach each node
	costSoFar := make(map[int]float32)
	costSoFar[startID] = 0

	// Map to reconstruct the path
	previousNode := make(map[int]int)

	for len(pq) > 0 {
		// Extract the node with the lowest cost
		current := pq[0]
		pq = pq[1:]

		// If we reached the goal, reconstruct the path as positions
		if current.nodeID == goalID {
			var path []pos
			for step := &current; step != nil; step = step.previous {
				node := findNodeByID(step.nodeID)
				if node != nil {
					path = append([]pos{node.pos}, path...)
				}
			}
			return path
		}

		// Explore neighbors
		for _, p := range game.currentmap.paths {
			var neighborID int
			if p.nodeA.id == current.nodeID {
				neighborID = p.nodeB.id
			} else if p.nodeB.id == current.nodeID {
				neighborID = p.nodeA.id
			} else {
				continue
			}

			// Calculate new cost
			newCost := costSoFar[current.nodeID] + p.cost
			if oldCost, exists := costSoFar[neighborID]; !exists || newCost < oldCost {
				costSoFar[neighborID] = newCost
				previousNode[neighborID] = current.nodeID

				// Add to priority queue
				pq = append(pq, priorityQueueItem{
					nodeID:   neighborID,
					cost:     newCost,
					previous: &current,
				})

				// Sort the priority queue by cost
				sort.Slice(pq, func(i, j int) bool {
					return pq[i].cost < pq[j].cost
				})
			}
		}
	}

	// If no path is found, return an empty slice
	return nil
}
