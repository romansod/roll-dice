/*
probgen.go - probability generator

Handles all probability event generation and
computation
*/
package probgen

import (
	"math/rand"
)

type ProbEvent struct {
	numEvents int
	outcomes  []string
	prng      func(int) int
}

func (pe ProbEvent) getProbValue() string {
	return pe.outcomes[pe.prng(len(pe.outcomes))]
}

func (pe ProbEvent) produceEvent(out chan string) {
	defer close(out)
	for i := 0; i < pe.numEvents; i++ {
		out <- pe.getProbValue()
	}
}

func (pe ProbEvent) consumeEvents(in chan string) map[string]int {
	results := make(map[string]int)

	for event := range in {
		results[event]++
	}

	return results
}

func (pe ProbEvent) computeProbability() map[string]int {
	events := make(chan string)

	go pe.produceEvent(events)

	return pe.consumeEvents(events)
}

func randNumGen(num_of_events int) int {
	return rand.Intn(num_of_events)
}

func GenerateProbabilisticEvent(events int, possibilities []string) map[string]int {
	probEvent := ProbEvent{numEvents: events, outcomes: possibilities, prng: randNumGen}
	return probEvent.computeProbability()
}
