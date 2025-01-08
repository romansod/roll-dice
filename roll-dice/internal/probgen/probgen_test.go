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
//
// NOTE: this must be run at the beginning of tests using
// PRNG_for_testing results
//
//	 Params
//		 rng_nums []int : slice of deterministic rng results
func initHardcodedRngNums(rng_nums []int) {
	hardcoded_rng_nums, hardcoded_rng_num_i = rng_nums, 0
}

// Retrieve the deterministic pregenerated rng results at the
// given index
//
//	 Params
//		 pe ProbEvent : the probability event related to the rng results
//		 index int    : the desired pregenerated rng result
//	 Returns
//		 string : specific event outcome
func getSpecificEvent(pe ProbEvent, index int) string {
	return pe.outcomes[hardcoded_rng_nums[index]%len(pe.outcomes)]
}

// Retrieve the next deterministic pregenerated rng result and
// advance the iterator
//
//	 Returns
//		 int : the current pregenerated rng result
func getHardcodedRngNum() int {
	next := hardcoded_rng_nums[hardcoded_rng_num_i]
	hardcoded_rng_num_i++
	return next
}

// The Psuedo Random Number Generator used for testing purposes
// and injected into ProbGen prng
func PRNG_for_testing(num_outcomes int) int {
	return getHardcodedRngNum() % num_outcomes
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
		outcomes: []string{Heads, Tails},
		prng:     PRNG_for_testing}

	// 0 -> 0
	expected, actual := 0, pe.prng(len(pe.outcomes))
	testing_utils.AssertEQi(t, expected, actual)

	// 3 -> 1
	expected, actual = 1, pe.prng(len(pe.outcomes))
	testing_utils.AssertEQi(t, expected, actual)

	// 5 -> 1
	expected, actual = 1, pe.prng(len(pe.outcomes))
	testing_utils.AssertEQi(t, expected, actual)

	// 22 -> 0
	expected, actual = 0, pe.prng(len(pe.outcomes))
	testing_utils.AssertEQi(t, expected, actual)

	// 7 -> 1
	expected, actual = 1, pe.prng(len(pe.outcomes))
	testing_utils.AssertEQi(t, expected, actual)

	// 4 -> 0
	expected, actual = 0, pe.prng(len(pe.outcomes))
	testing_utils.AssertEQi(t, expected, actual)
}

func TestGetProbValue(t *testing.T) {
	// This tests the probability to outcome conversion
	//
	// - 6 coin flip test

	initHardcodedRngNums([]int{0, 3, 5, 22, 7, 4})
	pe := ProbEvent{
		numEvents: 6,
		outcomes:  []string{Heads, Tails},
		prng:      PRNG_for_testing}

	// 0 -> heads
	expected, actual := Heads, pe.getProbValue()
	testing_utils.AssertEQ(t, expected, actual)

	// 3 -> tails
	expected, actual = Tails, pe.getProbValue()
	testing_utils.AssertEQ(t, expected, actual)

	// 5 -> tails
	expected, actual = Tails, pe.getProbValue()
	testing_utils.AssertEQ(t, expected, actual)

	// 22 -> heads
	expected, actual = Heads, pe.getProbValue()
	testing_utils.AssertEQ(t, expected, actual)

	// 7 -> tails
	expected, actual = Tails, pe.getProbValue()
	testing_utils.AssertEQ(t, expected, actual)

	// 4 -> heads
	expected, actual = Heads, pe.getProbValue()
	testing_utils.AssertEQ(t, expected, actual)
}

func TestProduceEvent(t *testing.T) {
	// This tests the production of probability events
	//
	// - 6 coin flip test

	initHardcodedRngNums([]int{0, 3, 5, 22, 7, 4})
	pe := ProbEvent{
		numEvents: 6,
		outcomes:  []string{Heads, Tails},
		prng:      PRNG_for_testing}

	events := make(chan string)

	go pe.produceEvent(events)

	num_e := 0

	for event := range events {
		expected, actual := getSpecificEvent(pe, num_e), event
		testing_utils.AssertEQ(t, expected, actual)

		num_e++
	}

	expected, actual := 6, num_e
	testing_utils.AssertEQi(t, expected, actual)
}

func TestConsumeEvent(t *testing.T) {
	// This tests the consumption of probability events
	//
	// - 6 coin flip test

	initHardcodedRngNums([]int{0, 3, 5, 22, 7, 4})
	pe := ProbEvent{
		numEvents: 6,
		outcomes:  []string{Heads, Tails},
		prng:      PRNG_for_testing}

	events := make(chan string)

	go pe.produceEvent(events)

	results := pe.consumeEvents(events)

	expected, actual := 3, results[Heads]
	testing_utils.AssertEQi(t, expected, actual)

	expected, actual = 3, results[Tails]
	testing_utils.AssertEQi(t, expected, actual)
}

func TestNegProbEventTypeValidate(t *testing.T) {
	// Tests that invalid input types are sufficiently handled

	// Invalid number of events (negative)
	ok, err := validate(CoinFlip{numEvents: -1})
	expected, actual := ErrInvalidEvents, err.Error()
	testing_utils.AssertEQb(t, false, ok)
	testing_utils.AssertEQ(t, expected, actual)

	ok, err = validate(DiceRoll{numEvents: -1, numSides: D4})
	expected, actual = ErrInvalidEvents, err.Error()
	testing_utils.AssertEQb(t, false, ok)
	testing_utils.AssertEQ(t, expected, actual)

	// Invalid number of events (zero)
	ok, err = validate(CoinFlip{numEvents: 0})
	expected, actual = ErrInvalidEvents, err.Error()
	testing_utils.AssertEQb(t, false, ok)
	testing_utils.AssertEQ(t, expected, actual)

	ok, err = validate(DiceRoll{numEvents: 0, numSides: D4})
	expected, actual = ErrInvalidEvents, err.Error()
	testing_utils.AssertEQb(t, false, ok)
	testing_utils.AssertEQ(t, expected, actual)
}

func TestPosProbEventTypeValidate(t *testing.T) {
	// Tests that valid inputs will not fail.
	//
	// NOTE: we do not check these results as they are not
	// deterministic. Please see TestGrnProbEvent* tests
	// which make the internal behavior deterministic
	// exclusively for testing

	ok, err := validate(CoinFlip{numEvents: 4})
	testing_utils.AssertEQb(t, true, ok)
	testing_utils.AssertNIL(t, err)

	ok, err = validate(DiceRoll{numEvents: 4, numSides: D6})
	testing_utils.AssertEQb(t, true, ok)
	testing_utils.AssertNIL(t, err)
}

func TestDiceRollValidate(t *testing.T) {
	// Test the validation of dice types

	// Invalid dice types

	// Smaller than all valid dice types
	diceRoll := NewDiceRoll(3, 3)
	ok, err := diceRoll.validate()
	testing_utils.AssertEQb(t, false, ok)
	testing_utils.AssertEQ(t, ErrInvalidDiceType, err.Error())

	// Between two valid dice types
	diceRoll = NewDiceRoll(3, 5)
	ok, err = diceRoll.validate()
	testing_utils.AssertEQb(t, false, ok)
	testing_utils.AssertEQ(t, ErrInvalidDiceType, err.Error())

	// Greater than the largest dice type
	diceRoll = NewDiceRoll(3, 21)
	ok, err = diceRoll.validate()
	testing_utils.AssertEQb(t, false, ok)
	testing_utils.AssertEQ(t, ErrInvalidDiceType, err.Error())

	// Valid dice types (D4, D6, D10, D12, D20)

	diceRoll = NewDiceRoll(3, D4)
	ok, _ = diceRoll.validate()
	testing_utils.AssertEQb(t, true, ok)

	diceRoll = NewDiceRoll(3, D6)
	ok, _ = diceRoll.validate()
	testing_utils.AssertEQb(t, true, ok)

	diceRoll = NewDiceRoll(3, D10)
	ok, _ = diceRoll.validate()
	testing_utils.AssertEQb(t, true, ok)

	diceRoll = NewDiceRoll(3, D12)
	ok, _ = diceRoll.validate()
	testing_utils.AssertEQb(t, true, ok)

	diceRoll = NewDiceRoll(3, D20)
	ok, _ = diceRoll.validate()
	testing_utils.AssertEQb(t, true, ok)
}

func TestGrnProbEventCoinFlip(t *testing.T) {
	// This tests the full production -> consumption of
	// probhen.ProbEvent.computeProbability
	//
	// - 6 coin flip test

	initHardcodedRngNums([]int{0, 3, 5, 22, 7, 4})
	coinFlip := ProbEvent{
		numEvents: 6,
		outcomes:  []string{Heads, Tails},
		prng:      PRNG_for_testing}

	res := coinFlip.computeProbability()

	// 3, 5, 7 -> 3 x tails
	expected, actual := 3, res[Tails]
	testing_utils.AssertEQi(t, expected, actual)

	// 1, 22, 4 -> 3 x tails
	expected, actual = 3, res[Heads]
	testing_utils.AssertEQi(t, expected, actual)
}

func TestGenProbEventDiceRoll(t *testing.T) {
	// This tests the full production -> consumption of
	// probhen.ProbEvent.computeProbability
	//
	// - 6 dice roll test

	initHardcodedRngNums([]int{0, 3, 5, 22, 7, 4})
	diceRoll := ProbEvent{
		numEvents: 6,
		outcomes:  []string{"1", "2", "3", "4", "5", "6"},
		prng:      PRNG_for_testing}

	res := diceRoll.computeProbability()

	// 0 -> 1 x 1
	expected, actual := 1, res["1"]
	testing_utils.AssertEQi(t, expected, actual)

	// 7 -> 1 x 2
	expected, actual = 1, res["2"]
	testing_utils.AssertEQi(t, expected, actual)

	// No roll equivalent to 3 -> 0
	expected, actual = 0, res["3"]
	testing_utils.AssertEQi(t, expected, actual)

	// 3 -> 1 x 4
	expected, actual = 1, res["4"]
	testing_utils.AssertEQi(t, expected, actual)

	// 6, 22 -> 2 x 5
	expected, actual = 2, res["5"]
	testing_utils.AssertEQi(t, expected, actual)

	// 5 -> 1 x 6
	expected, actual = 1, res["6"]
	testing_utils.AssertEQi(t, expected, actual)
}

func TestGenProbDisplaysCoinFlip(t *testing.T) {
	// Test the display functions of ProbEventTypes

	// CoinFlip
	origStdout, r, w := testing_utils.RedirectStdout()
	// 1) Flip just one coin, make sure percent is formatted
	coinFlip := CoinFlip{numEvents: 1}
	coinFlip.display(
		map[string]int{
			Heads: 1,
			Tails: 0,
		})

	output := testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	expected :=
		"(H) : 100.000000% : 1\n" +
			"(T) :   0.000000% : 0\n\n"
	testing_utils.AssertEQ(t, expected, output)

	// 2) Small scale should have round numbers
	origStdout, r, w = testing_utils.RedirectStdout()
	coinFlip = CoinFlip{numEvents: 10}
	coinFlip.display(
		map[string]int{
			Heads: 4,
			Tails: 6,
		})

	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	expected =
		"(H) :  40.000000% : 4\n" +
			"(T) :  60.000000% : 6\n\n"
	testing_utils.AssertEQ(t, expected, output)

	// 3) Large scale, non round should handle large values
	origStdout, r, w = testing_utils.RedirectStdout()
	coinFlip = CoinFlip{numEvents: 1000005}
	coinFlip.display(
		map[string]int{
			Heads: 499761,
			Tails: 500244,
		})

	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	expected =
		"(H) :  49.975849% : 499761\n" +
			"(T) :  50.024151% : 500244\n\n"
	testing_utils.AssertEQ(t, expected, output)
}

func TestGenProbDisplaysDiceRoll(t *testing.T) {
	// Test the display functions of ProbEventTypes

	// DiceRoll
	origStdout, r, w := testing_utils.RedirectStdout()
	// 1) Flip just one coin, make sure percent is formatted
	diceRoll := DiceRoll{numEvents: 1, numSides: D6}
	diceRoll.display(
		map[string]int{
			"1": 1,
			"2": 0,
			"3": 0,
			"4": 0,
			"5": 0,
			"6": 0,
		})

	output := testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	expected :=
		"[1]  : 100.000000% : 1\n" +
			"[2]  :   0.000000% : 0\n" +
			"[3]  :   0.000000% : 0\n" +
			"[4]  :   0.000000% : 0\n" +
			"[5]  :   0.000000% : 0\n" +
			"[6]  :   0.000000% : 0\n\n"
	testing_utils.AssertEQ(t, expected, output)

	// 2) Small scale should have round numbers
	origStdout, r, w = testing_utils.RedirectStdout()
	diceRoll = DiceRoll{numEvents: 10, numSides: D12}
	diceRoll.display(
		map[string]int{
			"1":  2,
			"2":  0,
			"3":  4,
			"4":  0,
			"5":  2,
			"6":  0,
			"7":  1,
			"8":  0,
			"9":  0,
			"10": 0,
			"11": 1,
			"12": 1,
		})

	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	expected =
		"[1]  :  20.000000% : 2\n" +
			"[2]  :   0.000000% : 0\n" +
			"[3]  :  40.000000% : 4\n" +
			"[4]  :   0.000000% : 0\n" +
			"[5]  :  20.000000% : 2\n" +
			"[6]  :   0.000000% : 0\n" +
			"[7]  :  10.000000% : 1\n" +
			"[8]  :   0.000000% : 0\n" +
			"[9]  :   0.000000% : 0\n" +
			"[10] :   0.000000% : 0\n" +
			"[11] :  10.000000% : 1\n" +
			"[12] :  10.000000% : 1\n\n"
	testing_utils.AssertEQ(t, expected, output)

	// 3) Large scale, non round should handle large values
	origStdout, r, w = testing_utils.RedirectStdout()
	diceRoll = DiceRoll{numEvents: 1000005, numSides: D4}
	diceRoll.display(
		map[string]int{
			"1": 250065,
			"2": 249826,
			"3": 249575,
			"4": 250539,
		})

	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	expected =
		"[1]  :  25.006374% : 250065\n" +
			"[2]  :  24.982475% : 249826\n" +
			"[3]  :  24.957375% : 249575\n" +
			"[4]  :  25.053774% : 250539\n\n"
	testing_utils.AssertEQ(t, expected, output)
}

func TestDisplayOneCoinFlip(t *testing.T) {
	// Test the proper coin handling for single action

	// Single coin (+)
	origStdout, r, w := testing_utils.RedirectStdout()

	DisplayOneFlipAction()
	output := testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	// Since this is non deterministic, just check if it is
	// one of the acceptable results
	testing_utils.AssertEQb(t, true, testing_utils.ContainsV(coinVisuals, output))
}

func TestDisplayOneDiceRoll(t *testing.T) {
	// Test the proper supported and unsupported dice type handling for single action

	// Single D4 (-)
	origStdout, r, w := testing_utils.RedirectStdout()

	ExecuteAndDisplayOneRollAction(D4)
	output := testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQb(t, false, output != ErrUnsupportedDiceType)

	// Single D6 (+)
	origStdout, r, w = testing_utils.RedirectStdout()

	ExecuteAndDisplayOneRollAction(D6)
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	// Since this is non deterministic, just check if it is
	// one of the acceptable results
	testing_utils.AssertEQb(t, true, testing_utils.ContainsV(d6Visuals, output))

	// Single D10 (-)
	origStdout, r, w = testing_utils.RedirectStdout()

	ExecuteAndDisplayOneRollAction(D10)
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQb(t, false, output != ErrUnsupportedDiceType)

	// Single D12 (-)
	origStdout, r, w = testing_utils.RedirectStdout()

	ExecuteAndDisplayOneRollAction(D12)
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQb(t, false, output != ErrUnsupportedDiceType)

	// Single D20 (-)
	origStdout, r, w = testing_utils.RedirectStdout()

	ExecuteAndDisplayOneRollAction(D20)
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQb(t, false, output != ErrUnsupportedDiceType)
}
