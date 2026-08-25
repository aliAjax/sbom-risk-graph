package domain

import "fmt"

type Graph struct {
	Nodes map[string]Component
	Edges map[string][]string
}

func NewGraph() *Graph {
	return &Graph{Nodes: make(map[string]Component), Edges: make(map[string][]string)}
}
func (g *Graph) AddComponent(c Component) error {
	if err := c.Validate(); err != nil {
		return err
	}
	g.Nodes[c.Key()] = c
	return nil
}
func (g *Graph) AddDependency(from, to string) error {
	if _, ok := g.Nodes[from]; !ok {
		return fmt.Errorf("source component %s missing", from)
	}
	if _, ok := g.Nodes[to]; !ok {
		return fmt.Errorf("target component %s missing", to)
	}
	g.Edges[from] = append(g.Edges[from], to)
	return nil
}
func (g *Graph) PathsFrom(root string, maxDepth int) [][]string {
	if maxDepth <= 0 {
		maxDepth = 32
	}
	var out [][]string
	var walk func(string, []string, map[string]bool)
	walk = func(node string, path []string, seen map[string]bool) {
		if len(path) > maxDepth {
			return
		}
		next := append(path, node)
		if seen[node] {
			return
		}
		seen[node] = true
		children := g.Edges[node]
		if len(children) == 0 {
			out = append(out, next)
			delete(seen, node)
			return
		}
		for _, child := range children {
			walk(child, next, seen)
		}
		delete(seen, node)
	}
	walk(root, nil, make(map[string]bool))
	return out
}
