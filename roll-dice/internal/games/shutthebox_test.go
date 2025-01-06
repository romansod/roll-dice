package games

import (
	"fmt"
	"testing"

	"github.com/romansod/roll-dice/internal/testing_utils"
)

// Helper function for this test file to batch test the existenec of a solution
// for a particular game state for every possible dice roll combination
//
//	Params
//		bitset int : game state integer. Generated by ConvertSlotsToGameState
//		expected [11] bool : splice indicating expected results slot solutions
//	Returns
//		error : any errors encountered including mismatch between expected
//		and actual results
func CheckPermutations(bitset int, expected [11]bool) error {
	// target : represents all values produced by 2 dice [2,12]
	// exp_i  : index through expected results

	for target, exp_i := 2, 0; target < 13; target++ {
		bitset_c := bitset
		exp, act := expected[exp_i], TargetSumExists(&bitset_c, target)
		if exp != act {
			return fmt.Errorf(
				"\ntarget = %d\nexp_i = %d\nexp   = %t\nact   = %t",
				target,
				exp_i,
				exp,
				act,
			)
		}

		exp_i++
	}

	return nil
}

func TestIsBitSet(t *testing.T) {
	// Bit checking tests

	testing_utils.AssertEQb(t, false, IsBitSet(OpenBox, 9))

	// Verify whether bits are set
	for i := 0; i < SizeBox; i++ {
		testing_utils.AssertEQb(t, true, IsBitSet(OpenBox, i))
	}

	// Verify whether bits are NOT set
	for i := 0; i < SizeBox; i++ {
		testing_utils.AssertEQb(t, false, IsBitSet(ShutBox, i))
	}
}

func TestGetSlotForPrint(t *testing.T) {
	// Slot retrieval for printing checks

	// Retrieve all slots (open)
	for i := 0; i < SizeBox; i++ {
		testing_utils.AssertEQ(
			t,
			"["+fmt.Sprintf("%d", i+1)+"]", // [1] -> [9]
			GetSlotForPrint(OpenBox, i))
	}

	// Retrive all slots (shut)
	for i := 0; i < SizeBox; i++ {
		testing_utils.AssertEQ(
			t,
			"["+EmptySlot+"]", // [_] x 9
			GetSlotForPrint(ShutBox, i))
	}
}

func TestSetBitEmpty(t *testing.T) {
	// Check the setting of bits to empty

	setvar := OpenBox
	// Start with completely open box
	testing_utils.AssertEQb(t, false, IsBoxEmpty(setvar))

	// Set each bit to empty one at a time
	for bit := 0; bit < SizeBox; bit++ {
		// Verify bit is set before
		testing_utils.AssertEQb(t, true, IsBitSet(setvar, bit))
		SetBitEmpty(&setvar, bit)
		// Verify bit is empty after
		testing_utils.AssertEQb(t, false, IsBitSet(setvar, bit))
	}

	// Verify the box is completely empty now
	testing_utils.AssertEQb(t, true, IsBoxEmpty(setvar))
}

func TestCombineDigits(t *testing.T) {
	// Test the helper function for combining digits from
	// a string. Inputs are also verified

	// (-) Invalid inputs
	output, err := combineDigits("")
	testing_utils.AssertEQ(t, ErrInvDigit, err.Error())
	testing_utils.AssertEQi(t, -1, output)

	output, err = combineDigits("1a345")
	testing_utils.AssertEQ(t, ErrInvDigit, err.Error())
	testing_utils.AssertEQi(t, -1, output)

	output, err = combineDigits("asdf")
	testing_utils.AssertEQ(t, ErrInvDigit, err.Error())
	testing_utils.AssertEQi(t, -1, output)

	output, err = combineDigits("-2")
	testing_utils.AssertEQ(t, ErrInvDigit, err.Error())
	testing_utils.AssertEQi(t, -1, output)

	output, err = combineDigits("0")
	testing_utils.AssertEQ(t, ErrInvDigit, err.Error())
	testing_utils.AssertEQi(t, -1, output)

	output, err = combineDigits("4209")
	testing_utils.AssertEQ(t, ErrInvDigit, err.Error())
	testing_utils.AssertEQi(t, -1, output)

	// (+) Valid inputs
	output, err = combineDigits("1")
	testing_utils.AssertNIL(t, err)
	testing_utils.AssertEQi(t, 1, output)

	output, err = combineDigits("145")
	testing_utils.AssertNIL(t, err)
	testing_utils.AssertEQi(t, 10, output)

	output, err = combineDigits("12345")
	testing_utils.AssertNIL(t, err)
	testing_utils.AssertEQi(t, 15, output)

	output, err = combineDigits("4597362")
	testing_utils.AssertNIL(t, err)
	testing_utils.AssertEQi(t, 36, output)
}

func TestIsValidShutInput(t *testing.T) {
	// Test the verification of input from a players turn

	// Ignoring stdout helps with extra lines added to processInput
	origStdout, ignoreOut := testing_utils.IgnoreStdout()

	// (-) Invalid inputs
	output := isValidShutInput("", 6)
	testing_utils.AssertEQb(t, false, output)

	output = isValidShutInput("1a345", 6)
	testing_utils.AssertEQb(t, false, output)

	output = isValidShutInput("asdf", 6)
	testing_utils.AssertEQb(t, false, output)

	output = isValidShutInput("-2", 6)
	testing_utils.AssertEQb(t, false, output)

	output = isValidShutInput("0", 6)
	testing_utils.AssertEQb(t, false, output)

	output = isValidShutInput("4209", 6)
	testing_utils.AssertEQb(t, false, output)

	// (-) Combined != Target
	output = isValidShutInput("1", 6)
	testing_utils.AssertEQb(t, false, output)

	output = isValidShutInput("145", 6)
	testing_utils.AssertEQb(t, false, output)

	output = isValidShutInput("12345", 6)
	testing_utils.AssertEQb(t, false, output)

	// (+) Combined == Target
	output = isValidShutInput("1", 1)
	testing_utils.AssertEQb(t, true, output)

	output = isValidShutInput("145", 10)
	testing_utils.AssertEQb(t, true, output)

	output = isValidShutInput("12345", 15)
	testing_utils.AssertEQb(t, true, output)

	output = isValidShutInput("4597362", 36)
	testing_utils.AssertEQb(t, true, output)

	testing_utils.IgnoreStdoutClose(origStdout, ignoreOut)
}

func TestUpdateGameState(t *testing.T) {
	// From an open box, update the game state until the
	// box is closed

	// Capture the print output for testing
	origStdout, r, w := testing_utils.RedirectStdout()
	stb := NewShutBox()
	stb.printGameState()
	output := testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "[1][2][3][4][5][6][7][8][9]", output)
	testing_utils.AssertEQb(t, false, stb.checkWinCondition())

	origStdout, r, w = testing_utils.RedirectStdout()
	stb.updateGameState("4")
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "[1][2][3][_][5][6][7][8][9]", output)
	testing_utils.AssertEQb(t, false, stb.checkWinCondition())

	origStdout, r, w = testing_utils.RedirectStdout()
	stb.updateGameState("17")
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "[_][2][3][_][5][6][_][8][9]", output)
	testing_utils.AssertEQb(t, false, stb.checkWinCondition())

	origStdout, r, w = testing_utils.RedirectStdout()
	stb.updateGameState("235")
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "[_][_][_][_][_][6][_][8][9]", output)
	testing_utils.AssertEQb(t, false, stb.checkWinCondition())

	origStdout, r, w = testing_utils.RedirectStdout()
	stb.updateGameState("9")
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "[_][_][_][_][_][6][_][8][_]", output)
	testing_utils.AssertEQb(t, false, stb.checkWinCondition())

	origStdout, r, w = testing_utils.RedirectStdout()
	stb.updateGameState("8")
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "[_][_][_][_][_][6][_][_][_]", output)
	testing_utils.AssertEQb(t, false, stb.checkWinCondition())

	origStdout, r, w = testing_utils.RedirectStdout()
	stb.updateGameState("6")
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "[_][_][_][_][_][_][_][_][_]", output)
	testing_utils.AssertEQb(t, true, stb.checkWinCondition())
}

func TestTargetSumExists(t *testing.T) {
	// Check that we can detect whether solutions exist even
	// with composite solutions. Last three slots will always
	// be composites

	// Fully open box
	bitset := ConvertSlotsToGameState("[1][2][3][4][5][6][7][8][9]")
	testing_utils.AssertNIL(
		t,
		CheckPermutations(
			bitset,
			[11]bool{
				true,
				true,
				true,
				true,
				true,
				true,
				true,
				true,
				true, // Composite
				true, // Composite
				true, // Composite
			},
		),
	)

	// Composite of 2 covers the missing slot 4
	bitset = ConvertSlotsToGameState("[1][2][3][_][5][6][7][8][9]")
	testing_utils.AssertNIL(
		t,
		CheckPermutations(
			bitset,
			[11]bool{
				true,
				true,
				true, // Composite
				true,
				true,
				true,
				true,
				true,
				true, // Composite
				true, // Composite
				true, // Composite
			},
		),
	)

	// Four 2 slot composites
	bitset = ConvertSlotsToGameState("[_][2][3][_][5][6][_][8][9]")
	testing_utils.AssertNIL(
		t,
		CheckPermutations(
			bitset,
			[11]bool{
				true,
				true,
				false,
				true,
				true,
				true, // Composite
				true,
				true,
				true, // Composite
				true, // Composite
				true, // Composite
			},
		),
	)

	// 3 slots
	bitset = ConvertSlotsToGameState("[_][_][_][_][_][6][_][8][9]")
	testing_utils.AssertNIL(
		t,
		CheckPermutations(
			bitset,
			[11]bool{
				false,
				false,
				false,
				false,
				true,
				false,
				true,
				true,
				false,
				false,
				false,
			},
		),
	)

	// Two slots
	bitset = ConvertSlotsToGameState("[_][_][_][_][_][6][_][8][_]")
	testing_utils.AssertNIL(
		t,
		CheckPermutations(
			bitset,
			[11]bool{
				false,
				false,
				false,
				false,
				true,
				false,
				true,
				false,
				false,
				false,
				false,
			},
		),
	)

	// Single slot
	bitset = ConvertSlotsToGameState("[_][_][_][_][_][6][_][_][_]")
	testing_utils.AssertNIL(
		t,
		CheckPermutations(
			bitset,
			[11]bool{
				false,
				false,
				false,
				false,
				true,
				false,
				false,
				false,
				false,
				false,
				false,
			},
		),
	)

	// Shut box
	bitset = ConvertSlotsToGameState("[_][_][_][_][_][_][_][_][_]")
	testing_utils.AssertNIL(
		t,
		CheckPermutations(
			bitset,
			[11]bool{
				false,
				false,
				false,
				false,
				false,
				false,
				false,
				false,
				false,
				false,
				false,
			},
		),
	)

	// Composite of 2 slots for 9
	bitset = ConvertSlotsToGameState("[1][_][_][_][_][_][_][8][_]")
	testing_utils.AssertNIL(
		t,
		CheckPermutations(
			bitset,
			[11]bool{
				false,
				false,
				false,
				false,
				false,
				false,
				true,
				true, // Composite
				false,
				false,
				false,
			},
		),
	)

	// Composite of 3 slots for 10
	bitset = ConvertSlotsToGameState("[_][2][3][_][5][_][_][_][_]")
	testing_utils.AssertNIL(
		t,
		CheckPermutations(
			bitset,
			[11]bool{
				true,
				true,
				false,
				true,
				false,
				true, // Composite
				true, // Composite
				false,
				true, // Composite
				false,
				false,
			},
		),
	)
}
