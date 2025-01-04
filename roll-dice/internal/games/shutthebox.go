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

		SetBitEmpty(&shutTheBox.gameState, digit_i-1)
	}
}

// Visualize the game state
func (shutTheBox ShutTheBox) printGameState() {
	gstate := ""
	for i := 0; i < SizeBox; i++ {
		gstate += GetSlotForPrint(shutTheBox.gameState, i)
	}

	fmt.Print(gstate)
}

// Check the win condition: box is shut
//
//	Returns
//		bool : true if box is totally shut, false otherwise
func (shutTheBox ShutTheBox) checkWinCondition() bool {
	return IsBoxEmpty(shutTheBox.gameState)
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
		slot_v = strconv.Itoa(slot + 1)
	}

	return fmt.Sprintf(Slot, slot_v)
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
