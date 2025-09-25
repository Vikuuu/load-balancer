package main

import "math/rand"

func weightedSelection(items []*Server) *Server {
	totalWeight := 0
	for _, item := range items {
		totalWeight += item.weight
	}

	r := rand.Intn(totalWeight)

	cursor := 0
	for _, item := range items {
		cursor += item.weight
		if cursor >= r {
			return item
		}
	}

	r = rand.Intn(len(items))
	return items[r]
}
