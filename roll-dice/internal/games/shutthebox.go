/*
shutthebox.go

ShutTheBox is a game which uses dice rolls
*/
package games

import (
	"errors"
	"fmt"
	"strconv"
)

/// Errors

const ErrInvDigit string = "invalid digit input not in range [1,9]"
const ErrNotEqTarget string = "input '%d' does not add up to target '%d'"

// Total number of slots
const SizeBox int = 9

// Initial open box
const OpenBox int = (1 << SizeBox) - 1

// Fully shut box
const ShutBox int = 0

// Slot display for formatting
const Slot string = "[%s]"

// Empty slot display value
const EmptySlot string = "_"

type ShutTheBox struct {
	gameState int
}

// Initialize private fields
//
//	Returns
//		*ShutTheBox : new ShutTheBox object
func NewShutBox() *ShutTheBox {
	return &ShutTheBox{
		gameState: OpenBox, // game state stored as 9 bits
	}
}

// Verified update to the game state in a single string
//
//	Params
//		update string : digits referring to the slots to be closed. Ex: "142"
func (shutTheBox *ShutTheBox) updateGameState(update string) {
	for _, d := range update {
		digit_i, _ := strconv.Atoi(string(d))

		SetBitEmpty(&shutTheBox.gameState, GetValueSlot(digit_i))
	}
}

// Visualize the game state
//
// Ex: slot 1, 4, 7 are closed:
//
//	[_][2][3][_][5][6][_][8][9]
func (shutTheBox ShutTheBox) printGameState() {
	fmt.Print(AssembleSlotsToDisplay(shutTheBox.gameState))
}

// Check the win condition: box is shut
//
//	Returns
//		bool : true if box is totally shut, false otherwise
func (shutTheBox ShutTheBox) checkWinCondition() bool {
	return IsBoxEmpty(shutTheBox.gameState)
}

// Check whether a solution exists
//
//	Params
//		target int : the target sum to check for among the open slots
//	Returns
//		bool : whether there is a solution in the current game state
//		to satisfy the target
func (shutTheBox ShutTheBox) checkSolutionExists(target int) bool {
	gstate := shutTheBox.gameState
	return TargetSumExists(&gstate, target)
}

// Check whether provided input is valid to satisfy the target
//
//	Params
//		shutInput string : digits referring to the slots to be closed. Ex: "142"
//		target int       : the target for the sum of the shutInput
//	Returns
//		bool : true if shutInput digits sum to the target
func isValidShutInput(shutInput string, target int) bool {
	// Sum the digits
	combinedDigits, err := combineDigits(shutInput)

	if err != nil {
		fmt.Print(err.Error())
		return false
	}

	// Verify whether the inputs actually add up to the target
	if combinedDigits != target {
		fmt.Printf(ErrNotEqTarget, combinedDigits, target)
		return false
	}

	return true
}

// Sum the individual digits. There are only valid inputs from [1,9]
//
// This function also checks for invalid inputs and returns errors
// for empty inputs and string to int conversions
//
//	Params
//		digits string : digits referring to the slots to be closed. Ex: "142"
//	Returns
//		int : the sum of the digits or -1 if error is encountered
//		error : any errors encountered or nil
func combineDigits(digits string) (int, error) {
	// Empty input string is invalid
	if digits == "" {
		return -1, errors.New(ErrInvDigit)
	}

	combinedDigits := 0

	// Add all the digits together
	for _, d := range digits {
		digit_i, err := strconv.Atoi(string(d))
		// Any error in the conversion or an invalid digit will
		// cause immediate termination of execution
		if err != nil || digit_i < 1 {
			return -1, errors.New(ErrInvDigit)
		}

		combinedDigits += digit_i
	}

	return combinedDigits, nil
}

// Check whether the bit in the bitset is on
//
//	Params
//		bitset int : bits to check
//		bit int    : the bit in bitset to check
//	Returns
//		bool : true if bit is on, false otherwise
func IsBitSet(bitset int, bit int) bool {
	return bitset&(1<<bit) != 0
}

// Set the given bit in the bitset to off
//
//	Params
//		bitset int : bits that contain the bit to set off
//		bit int    : the bit in bitset to turn off
func SetBitEmpty(bitset *int, bit int) {
	*bitset = *bitset &^ (1 << bit)
}

// Retrieve the visualized slot for printing
//
// Ex: Open slot   -> [0][1][2] ... [9]
// Ex: Closed slot -> [_]
//
//	Params
//		gstate int : game state bitset
//		slot int   : the slot we want to visualize
//	Returns
//		string : the visualized slot
func GetSlotForPrint(gstate int, slot int) string {
	slot_v := EmptySlot

	if IsBitSet(gstate, slot) {
		slot_v = strconv.Itoa(GetSlotValue(slot))
	}

	return fmt.Sprintf(Slot, slot_v)
}

// Get the value for the given slot index in the game state
//
//	Params
//		slot int : the slot in the game state who's value we want
//	Returns
//		int : the value of the requested slot
func GetSlotValue(slot int) int {
	return slot + 1
}

// Get the slot index for the given value in the game state
//
//	Params
//		value int : the value in the game state who's slot index we want
//	Returns
//		int : the slot index of the given value
func GetValueSlot(value int) int {
	return value - 1
}

// Create formatted display for the provided game state
//
//	 Ex: gstate(32) -> "[_][_][_][_][_][6][_][_][_]"
//		Params
//			gstate int : game state to display
//		Returns
//			string : display string
func AssembleSlotsToDisplay(gstate int) string {
	gstateslots := ""
	for i := 0; i < SizeBox; i++ {
		gstateslots += GetSlotForPrint(gstate, i)
	}

	return gstateslots
}

// Helper function to convert displayed game state to internal game state
//
// Useful in tests. Example: "[_][_][_][_][_][6][_][_][_]" -> 32
//
//	Params
//		gslots string : formatted slot display for conversion
//	Returns
//		int : game state representation
func ConvertSlotsToGameState(gslots string) int {
	gstate := 0
	for i := 0; i < SizeBox; i++ {
		// Turn each bit on for each open slot
		gstate |= (ConvertSlotToBit(gslots, i) << i)
	}

	return gstate
}

// Helper function to convertdisplayed slot to a bit
//
// Ex: [_] -> 0
// Ex: [4] -> 1
// Ex: [7] -> 1
//
//	Params
//		gslots string : game state as visual string (9 slots)
//		slot int      : the slot we want to convert to a bit (off or on)
//	Returns
//		int : 0 if [_] and 1 if [1->9]
func ConvertSlotToBit(gslots string, slot int) int {
	// [X][.]...
	// \_\
	//    \
	//     slot index 0 in the display format covers gslots[0:3]
	// We read each slot 3 characters at a time
	if gslots[slot*3:(slot*3)+3] == fmt.Sprintf(Slot, EmptySlot) {
		return 0
	} else {
		return 1
	}
}

// Recursive algorithm to find the first available potential solution in the
// provided bitset that satisfies the indicated target
//
// Terminal State : when target is satisfied by a given slot exactly
//
// Recursive Principle : split the target into a low and high that sum to the
// original target and recurse on those sub targets
//
// Game Rules : A slot can only be used once when trying to reach the target
//
//	Params
//		bitset *int : persistent game state used to reach target
//		target int  : the target sum of open slots in the game state
//	Returns
//		bool : true if a solution exists, false otherwise
func TargetSumExists(bitset *int, target int) bool {
	// Terminal State -> check if satisfied by single slot
	if IsBitSet(*bitset, GetValueSlot(target)) {
		// Consume this bit to satisfy the Game Rules. This
		// slot can not be used in lateral sub targets
		SetBitEmpty(bitset, GetValueSlot(target))
		return true
	}

	// We copy the incoming game state in case we need to reset due to
	// an unsolvable sub target
	orig_bitset := *bitset

	// Create 2 sum programmatically
	low_v := 1
	high_v := target - 1

	// Bisecting from the middle means we can avoid symmetric duplication
	for low_v < high_v {
		// Recurse over sub targets. If the first sub target does not exist
		// we can short circuit
		if TargetSumExists(bitset, low_v) && TargetSumExists(bitset, high_v) {
			return true
		} else {
			// Reset for the next iteration
			bitset = &orig_bitset
			// initialize the next permutation
			low_v++
			high_v--
		}
	}

	// No solution for the current target
	return false
}

// Check whether the box is empty, ie all slots closed
//
//	Params
//		gstate int : game state bitset
//	Returns
//		bool : true if all slots in box are shut
func IsBoxEmpty(gstate int) bool {
	return gstate == ShutBox
}
