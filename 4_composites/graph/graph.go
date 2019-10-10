// Package graph conceptualizes map usage to define a graph.
package graph

var graph = make(map[string]map[string]bool)

// idiomatic way to populate a map lazily
func addEdge(from, to string) {
	edges := graph[from]
	if edges == nil {
		edges = make(map[string]bool)
		graph[from] = edges
	}
	edges[to] = true
}

// shows usefulness of zero value of unitialized map entry
// always returns a meaningful result
func hasEdge(from, to string) bool {
	return graph[from][to]
}
