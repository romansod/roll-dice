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

const ErrInvalidDiceType = "invalid dice type: must be one of " + ValidDiceTypes

const (
	D4  = 4
	D6  = 6
	D10 = 10
	D12 = 12
	D20 = 20
)

const ValidDiceTypes = "(4, 6, 10, 12, 20)"

var DicePossibleValues = []string{
	"1", "2", "3", "4",
	"5", "6",
	"7", "8", "9", "10",
	"11", "12",
	"13", "14", "15", "16", "17", "18", "19", "20",
}

type DiceRoll struct {
	NumEvents int // number of coin flips
	NumSides  int // number of sides on dice
}

func (diceRoll DiceRoll) validate() (bool, error) {
	//  Need to make sure the provided dice type is valid
	if !validDiceType(diceRoll.NumSides) {
		return false, errors.New(ErrInvalidDiceType)
	}

	return true, nil
}

func (diceRoll DiceRoll) execute() error {
	res, err := GenerateProbabilisticEvent(
		diceRoll.NumEvents,
		possibleDiceValues(diceRoll.NumSides))

	if err == nil {
		diceRoll.display(res)
	}

	return err
}

// Print the coin flip results. Example:
//
// NumEvents: 123435
//
// (H) :  49.882126% : 61572
//
// (T) :  50.117874% : 61863
//
//	Params
//		res map[string]int : results of coin flips
func (diceRoll DiceRoll) display(res map[string]int) {
	for i := 1; i <= diceRoll.NumSides; i++ {
		i_s := strconv.Itoa(i)
		fmt.Printf(
			"%-4s : %10f%% : %d\n",
			"["+i_s+"]",
			Percent(res[i_s], diceRoll.NumEvents),
			res[i_s],
		)
	}
}

func validDiceType(dType int) bool {
	switch dType {
	case D4, D6, D10, D12, D20:
		return true
	default:
		return false
	}
}

func possibleDiceValues(dType int) []string {
	return DicePossibleValues[:dType]
}
