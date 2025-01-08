/*
diceroll.go

DiceRoll is a ProbEventType which
describes dice rolls
*/
package probgen

import (
	"errors"
	"fmt"
	"strconv"
)

const ErrInvalidDiceType = "invalid number of dice sides: must be one of " + ValidDiceTypes
const ErrUnsupportedDiceType = "unsupported dice type, only support D6 for now"

// Potential dice types
const (
	D4  = 4
	D6  = 6
	D10 = 10
	D12 = 12
	D20 = 20
)

// Valid dice types string
const ValidDiceTypes = "(4, 6, 10, 12, 20)"

// All possible dice values
var dicePossibleValues = []string{
	"1", "2", "3", "4", // D4
	"5", "6", // -> D6
	"7", "8", "9", "10", // -> D10
	"11", "12", // -> D12
	"13", "14", "15", "16", "17", "18", "19", "20", // -> D20
}

const (
	r1  = iota // 0
	r2         // 1
	r3         // 2
	r4         // 3
	r5         // 4
	r6         // 5
	r7         // 6
	r8         // 7
	r9         // 8
	r10        // 9
	r11        // 10
	r12        // 11
	r13        // 12
	r14        // 13
	r15        // 14
	r16        // 15
	r17        // 16
	r18        // 17
	r19        // 18
	r20        // 19
)

// All visual representations of 6 sided dice
var d6Visuals = map[int]string{
	r1: " -------\n" +
		"|       |\n" +
		"|   o   |\n" +
		"|       |\n" +
		" -------\n",
	r2: " -------\n" +
		"| o     |\n" +
		"|       |\n" +
		"|     o |\n" +
		" -------\n",
	r3: " -------\n" +
		"| o     |\n" +
		"|   o   |\n" +
		"|     o |\n" +
		" -------\n",
	r4: " -------\n" +
		"| o   o |\n" +
		"|       |\n" +
		"| o   o |\n" +
		" -------\n",
	r5: " -------\n" +
		"| o   o |\n" +
		"|   o   |\n" +
		"| o   o |\n" +
		" -------\n",
	r6: " -------\n" +
		"| o   o |\n" +
		"| o   o |\n" +
		"| o   o |\n" +
		" -------\n",
}

type DiceRoll struct {
	numEvents int // number of coin flips
	numSides  int // number of sides on dice
}

// Initialize private fields
//
//	Params
//		nEvents int : number of DiceRoll events
//		nSides int  : number of sides to the dice
//	Returns
//		*DiceRoll : new DiceRoll object
func NewDiceRoll(nEvents int, nSides int) *DiceRoll {
	return &DiceRoll{
		numEvents: nEvents,
		numSides:  nSides,
	}
}

func (diceRoll DiceRoll) validate() (bool, error) {
	//  Need to make sure the provided dice type is valid
	if !validDiceType(diceRoll.numSides) {
		return false, errors.New(ErrInvalidDiceType)
	}

	return true, nil
}

func (diceRoll DiceRoll) execute() error {
	res, err := GenerateProbabilisticEvent(
		diceRoll.numEvents,
		possibleDiceValues(diceRoll.numSides))

	if err == nil {
		diceRoll.display(res)
	}

	return err
}

// Exposed endpoint to execute one dice roll and
// print out a visual of the result
//
// NOTE: not all dice types are supported yet
//
//	Params
//		nSides int : indicate the number of sides for the dice
//	Returns
//		int : result of dice roll
func ExecuteAndDisplayOneRollAction(nSides int) int {
	res := ExecuteOneRollAction(nSides)

	// Only support D6 for now
	switch nSides {
	case D6:
		fmt.Print(d6Visuals[res])
		return res
	default:
		fmt.Print(ErrUnsupportedDiceType)
		return -1
	}
}

// One dice roll action
//
//	Params
//		nSides int : number of sides for the die
//	Returns
//		int : dice value 0 -> nSides - 1
func ExecuteOneRollAction(nSides int) int {
	pe := ProbEvent{
		numEvents: 1,
		outcomes:  possibleDiceValues(nSides),
		prng:      randNumGen}

	return pe.getProbValue()
}

// Print the dice roll results. Example:
//
// numEvents: 2
//
// numSides: 4
//
// [1]  :  50.00000% : 1
//
// [2] :   50.00000% : 1
//
// [3] :    0.00000% : 0
//
// [4] :    0.00000% : 0
//
//	Params
//		res map[string]int : results of dice rolls
func (diceRoll DiceRoll) display(res map[string]int) {
	for i := 1; i <= diceRoll.numSides; i++ {
		i_s := strconv.Itoa(i)
		fmt.Printf(
			"%-4s : %10f%% : %d\n",
			"["+i_s+"]",
			Percent(res[i_s], diceRoll.numEvents),
			res[i_s],
		)
	}
	fmt.Print("\n")
}

// Retrieve number of events
//
//	Returns
//		int : number of events
func (diceRoll DiceRoll) getNumEvents() int {
	return diceRoll.numEvents
}

// Make sure this is a valid dice type
//
//	Params
//		dType int : dice type
//	Returns
//		bool : true if dType in {D4, D6, D10, D12, D20}
//			   false otherwise
func validDiceType(dType int) bool {
	switch dType {
	case D4, D6, D10, D12, D20:
		return true
	default:
		return false
	}
}

// Get all the possible values for the particular type of dice. Do not call
// without validating first
//
//	Params
//		dType int : dice type
//	Returns
//		[]string : slice of valid dice outputs
//			Ex:
//				D4  [1, 4]
//				D6  [1, 6]
//				D10 [1, 10]
//				D12 [1, 12]
//				D20 [1, 20]
func possibleDiceValues(dType int) []string {
	return dicePossibleValues[:dType]
}
