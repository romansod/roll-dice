package probgen

import (
	"testing"

	"github.com/romansod/roll-dice/internal/testing_utils"
)

/// PRNG for testing
///
/// Use these global variables and functions for managing
/// a deterministic collection of hardcoded values to be
/// injected in place of the random number generator

// predetermined PRNG results
var hardcoded_rng_nums []int

// iterator through hardcoded_rng_nums
var hardcoded_rng_num_i int = 0

// Seed the deterministic slice of pregenerated rng results
func initHardcodedRngNums(rng_nums []int) {
	hardcoded_rng_nums, hardcoded_rng_num_i = rng_nums, 0
}

func getSpecificEvent(pe ProbEvent, index int) string {
	return pe.outcomes[hardcoded_rng_nums[index]%len(pe.outcomes)]
}

func getNextHardcodedRngNum() int {
	next := hardcoded_rng_nums[hardcoded_rng_num_i]
	hardcoded_rng_num_i++
	return next
}

func PRNG_for_testing(num_outcomes int) int {
	return getNextHardcodedRngNum() % num_outcomes
}

/// Tests for probgen

func TestPRNG_for_testing(t *testing.T) {
	// This test demonstrates the use of the PRNG_for_testing
	// behavior for simulatin a series of deterministic rng events
	//
	// - 6 coin flip test

	initHardcodedRngNums([]int{0, 3, 5, 22, 7, 4})
	pe := ProbEvent{
		//numEvents: 6, NOT USED
		outcomes: []string{"heads", "tails"},
		prng:     PRNG_for_testing}

	// 0 -> 0
	expected, actual := 0, pe.prng(len(pe.outcomes))
	if !testing_utils.AssertEQi(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}

	// 3 -> 1
	expected, actual = 1, pe.prng(len(pe.outcomes))
	if !testing_utils.AssertEQi(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}

	// 5 -> 1
	expected, actual = 1, pe.prng(len(pe.outcomes))
	if !testing_utils.AssertEQi(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}

	// 22 -> 0
	expected, actual = 0, pe.prng(len(pe.outcomes))
	if !testing_utils.AssertEQi(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}

	// 7 -> 1
	expected, actual = 1, pe.prng(len(pe.outcomes))
	if !testing_utils.AssertEQi(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}

	// 4 -> 0
	expected, actual = 0, pe.prng(len(pe.outcomes))
	if !testing_utils.AssertEQi(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}
}

func TestGetProbValue(t *testing.T) {
	// This tests the probability to outcome conversion
	//
	// - 6 coin flip test

	initHardcodedRngNums([]int{0, 3, 5, 22, 7, 4})
	pe := ProbEvent{
		numEvents: 6,
		outcomes:  []string{"heads", "tails"},
		prng:      PRNG_for_testing}

	// 0 -> heads
	expected, actual := "heads", pe.getProbValue()
	if !testing_utils.AssertEQ(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}

	// 3 -> tails
	expected, actual = "tails", pe.getProbValue()
	if !testing_utils.AssertEQ(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}

	// 5 -> tails
	expected, actual = "tails", pe.getProbValue()
	if !testing_utils.AssertEQ(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}

	// 22 -> heads
	expected, actual = "heads", pe.getProbValue()
	if !testing_utils.AssertEQ(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}

	// 7 -> tails
	expected, actual = "tails", pe.getProbValue()
	if !testing_utils.AssertEQ(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}

	// 4 -> heads
	expected, actual = "heads", pe.getProbValue()
	if !testing_utils.AssertEQ(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}
}

func TestProduceEvent(t *testing.T) {
	// This tests the production of probability events
	//
	// - 6 coin flip test

	initHardcodedRngNums([]int{0, 3, 5, 22, 7, 4})
	pe := ProbEvent{
		numEvents: 6,
		outcomes:  []string{"heads", "tails"},
		prng:      PRNG_for_testing}

	events := make(chan string)

	go pe.produceEvent(events)

	num_e := 0

	for event := range events {
		expected, actual := getSpecificEvent(pe, num_e), event
		if !testing_utils.AssertEQ(expected, actual) {
			t.Errorf(testing_utils.AssertFailed, expected, actual)
		}

		num_e++
	}

	expected, actual := 6, num_e
	if !testing_utils.AssertEQi(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}
}

func TestConsumeEvent(t *testing.T) {
	// This tests the consumption of probability events
	//
	// - 6 coin flip test

	initHardcodedRngNums([]int{0, 3, 5, 22, 7, 4})
	pe := ProbEvent{
		numEvents: 6,
		outcomes:  []string{"heads", "tails"},
		prng:      PRNG_for_testing}

	events := make(chan string)

	go pe.produceEvent(events)

	results := pe.consumeEvents(events)

	expected, actual := 3, results["heads"]
	if !testing_utils.AssertEQi(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}

	expected, actual = 3, results["tails"]
	if !testing_utils.AssertEQi(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}
}

func TestGrnProbEventCoinFlip(t *testing.T) {
	// This tests the full production -> consumption of
	// probhen.ProbEvent.computeProbability
	//
	// - 6 coin flip test

	initHardcodedRngNums([]int{0, 3, 5, 22, 7, 4})
	coinFlip := ProbEvent{
		numEvents: 6,
		outcomes:  []string{"heads", "tails"},
		prng:      PRNG_for_testing}

	res := coinFlip.computeProbability()

	// 3, 5, 7 -> 3 x tails
	expected, actual := 3, res["tails"]
	if !testing_utils.AssertEQi(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}

	// 1, 22, 4 -> 3 x tails
	expected, actual = 3, res["heads"]
	if !testing_utils.AssertEQi(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}
}

func TestGenProbEventRollDice(t *testing.T) {
	// This tests the full production -> consumption of
	// probhen.ProbEvent.computeProbability
	//
	// - 6 dice roll test

	initHardcodedRngNums([]int{0, 3, 5, 22, 7, 4})
	rollDice := ProbEvent{
		numEvents: 6,
		outcomes:  []string{"1", "2", "3", "4", "5", "6"},
		prng:      PRNG_for_testing}

	res := rollDice.computeProbability()

	// 0 -> 1 x 1
	expected, actual := 1, res["1"]
	if !testing_utils.AssertEQi(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}

	// 7 -> 1 x 2
	expected, actual = 1, res["2"]
	if !testing_utils.AssertEQi(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}

	// No roll equivalent to 3 -> 0
	expected, actual = 0, res["3"]
	if !testing_utils.AssertEQi(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}

	// 3 -> 1 x 4
	expected, actual = 1, res["4"]
	if !testing_utils.AssertEQi(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}

	// 6, 22 -> 2 x 5
	expected, actual = 2, res["5"]
	if !testing_utils.AssertEQi(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}

	// 5 -> 1 x 6
	expected, actual = 1, res["6"]
	if !testing_utils.AssertEQi(expected, actual) {
		t.Errorf(testing_utils.AssertFailed, expected, actual)
	}
}
