/*
probgen.go - probability generator

Handles all probability event generation and
computation
*/
package probgen

import (
	"errors"
	"math/rand"
)

/// Constants

const ErrInvalidEvents = "invalid number of events: must be more than one event"
const ErrInvalidPossibilities = "invalid number of possibilities: must have at least one possible outcome"

// Generic probability event object
type ProbEvent struct {
	numEvents int           // Number of probabilistic events
	outcomes  []string      // Total possible outcomes of events
	prng      func(int) int // The Pseudo Random Number Generator to use
}

// Get a probability value based on a random number generator bounded
// by the number of outcomes
//
//	Returns
//		int : a randomly selected number from [0, num_pe_outcomes]
func (pe ProbEvent) getProbValue() int {
	rngn := pe.prng(len(pe.outcomes))
	return rngn
}

func (pe ProbEvent) getProbOutcome(rngn int) string {
	return pe.outcomes[rngn]
}

// Add a probability computation to the out channel
//
//	Params
//		out chan string : output channel for probability computation results
func (pe ProbEvent) produceEvent(out chan string) {
	defer close(out)
	for i := 0; i < pe.numEvents; i++ {
		out <- pe.getProbOutcome(pe.getProbValue())
	}
}

// Process a probability computation result and aggregate based on that outcome
//
//	Params
//		in chan string : input channel containing all probability computations
//	Returns
//		map[string]int : aggregated results of the input channel
func (pe ProbEvent) consumeEvents(in chan string) map[string]int {
	results := make(map[string]int)

	for event := range in {
		results[event]++
	}

	return results
}

// Compute the probability for a ProbEvent based on its
// numEvents, outcomes, and prng()
//
//	Returns
//		map[string]int : aggregation of results into a
//		table that is indexed by the possible outcomes
//		and returns the number of times that outcome
//		occurred
func (pe ProbEvent) computeProbability() map[string]int {
	events := make(chan string)

	go pe.produceEvent(events)

	return pe.consumeEvents(events)
}

// Generate a random number bounded by the number of outcomes
//
// NOTE: num_outcomes is not zero based, but the possible outcomes
// are and this is handled by the half open interval: [0, n)
//
//	Returns
//		int : a number in the range: [0, n)
func randNumGen(num_outcomes int) int {
	return rand.Intn(num_outcomes)
}

// Given the number of events and the possible outcomes of the events, return
// a table of results
//
//	Params
//		events int : number of probability events taking place
//		possibilities []string : all the possible outcomes
//	Returns
//		map[string]int : aggregation of results into a
//		table that is indexed by the possible outcomes
//		and returns the number of times that outcome
//		occurred
//		error          : any errors encountered
//
//	Ex:
//		events       : 3
//		possibilities: {"heads", "tails"}
//
//		e1 : "heads"
//		e2 : "tails"
//		e3 : "heads"
//
//		returns : {"heads":2, "tails":1}
func GenerateProbabilisticEvent(events int, possibilities []string) (map[string]int, error) {
	if len(possibilities) < 1 {
		// Must have at least one possible outcome
		return nil, errors.New(ErrInvalidPossibilities)
	}

	probEvent := ProbEvent{numEvents: events, outcomes: possibilities, prng: randNumGen}
	return probEvent.computeProbability(), nil
}

// ProbEvent interface for use in options
type ProbEventType interface {
	validate() (bool, error) // Check input is valid
	execute() error          // Compute and display result
	display(map[string]int)  // Display results
	getNumEvents() int       // Retrieve number of events
}

func ValidateAndExecute(probEventType ProbEventType) error {
	// Generic probability event validation
	ok, err := validate(probEventType)
	if !ok {
		return err
	}

	// Specialized probability event validation
	ok, err = probEventType.validate()
	if !ok {
		return err
	}

	return probEventType.execute()
}

// Generally applicable ProbEvent validation
//
//	Params
//		probEventType ProbEventType : probability event to check
//	Returns
//		bool  : valid status
//		error : indicates any errors leading to validation failure
func validate(probEventType ProbEventType) (bool, error) {
	if probEventType.getNumEvents() < 1 {
		return false, errors.New(ErrInvalidEvents)
	}

	return true, nil
}

// Utility to compute the percent: numerator / denominator
//
//	Params
//		numerator int   : divided by denominator
//		denominator int : divides numerator
func Percent(numerator int, denominator int) float32 {
	return float32(numerator) * 100 / float32(denominator)
}
