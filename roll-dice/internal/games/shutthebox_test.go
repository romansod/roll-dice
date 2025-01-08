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

func TestIsValidShutInput(t *testing.T) {
	// Test the verification and processing of input from a players turn

	// Ignoring stdout helps with extra lines added to processInput
	origStdout, ignoreOut := testing_utils.IgnoreStdout()

	gstate := OpenBox

	// (-) Invalid inputs
	_, err := processProposedUpdate(gstate, "", 6)
	testing_utils.AssertEQ(t, ErrInvDigit, err.Error())

	gstate = OpenBox
	_, err = processProposedUpdate(gstate, "1a345", 6)
	testing_utils.AssertEQ(t, ErrInvDigit, err.Error())

	gstate = OpenBox
	_, err = processProposedUpdate(gstate, "asdf", 6)
	testing_utils.AssertEQ(t, ErrInvDigit, err.Error())

	gstate = OpenBox
	_, err = processProposedUpdate(gstate, "-2", 6)
	testing_utils.AssertEQ(t, ErrInvDigit, err.Error())

	gstate = OpenBox
	_, err = processProposedUpdate(gstate, "0", 6)
	testing_utils.AssertEQ(t, ErrInvDigit, err.Error())

	gstate = OpenBox
	_, err = processProposedUpdate(gstate, "4209", 6)
	testing_utils.AssertEQ(t, ErrInvDigit, err.Error())

	// (-) Combined != Target
	gstate = OpenBox
	_, err = processProposedUpdate(gstate, "1", 6)
	testing_utils.AssertEQ(t, fmt.Sprintf(ErrNotEqTarget, 1, 6), err.Error())

	gstate = OpenBox
	_, err = processProposedUpdate(gstate, "145", 6)
	testing_utils.AssertEQ(t, fmt.Sprintf(ErrNotEqTarget, 10, 6), err.Error())

	gstate = OpenBox
	_, err = processProposedUpdate(gstate, "12345", 6)
	testing_utils.AssertEQ(t, fmt.Sprintf(ErrNotEqTarget, 15, 6), err.Error())

	// (+) Combined == Target
	gstate = OpenBox
	gstate_processed, err := processProposedUpdate(gstate, "1", 1)
	testing_utils.AssertNIL(t, err)
	testing_utils.AssertEQ(t, "[_][2][3][4][5][6][7][8][9]", AssembleSlotsToDisplay(gstate_processed))

	gstate_processed, err = processProposedUpdate(gstate, "45", 9)
	testing_utils.AssertNIL(t, err)
	testing_utils.AssertEQ(t, "[1][2][3][_][_][6][7][8][9]", AssembleSlotsToDisplay(gstate_processed))

	gstate_processed, err = processProposedUpdate(gstate, "1245", 12)
	testing_utils.AssertNIL(t, err)
	testing_utils.AssertEQ(t, "[_][_][3][_][_][6][7][8][9]", AssembleSlotsToDisplay(gstate_processed))

	gstate_processed, err = processProposedUpdate(gstate, "134", 8)
	testing_utils.AssertNIL(t, err)
	testing_utils.AssertEQ(t, "[_][2][_][_][5][6][7][8][9]", AssembleSlotsToDisplay(gstate_processed))

	testing_utils.IgnoreStdoutClose(origStdout, ignoreOut)
}

func TestUpdateGameState(t *testing.T) {
	// From an open box, update the game state until the
	// box is closed. Also demonstrate no-op when update
	// fails

	// Capture the print output for testing
	origStdout, r, w := testing_utils.RedirectStdout()
	stb := NewShutBox([]string{"p1"})
	stb.printGameState()
	output := testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p1\n\n[1][2][3][4][5][6][7][8][9]", output)
	testing_utils.AssertEQb(t, false, stb.checkWinCondition())

	// (+) Successful update of single slot
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.updateGameState("4", 4)
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p1\n\n[1][2][3][_][5][6][7][8][9]", output)
	testing_utils.AssertEQb(t, false, stb.checkWinCondition())

	// (-) No op when given an update that is not valid
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.updateGameState("adg", 4)
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p1\n\n[1][2][3][_][5][6][7][8][9]", output)
	testing_utils.AssertEQb(t, false, stb.checkWinCondition())

	// (-) No op when given an update that does not satisfy the target
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.updateGameState("178", 4)
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p1\n\n[1][2][3][_][5][6][7][8][9]", output)
	testing_utils.AssertEQb(t, false, stb.checkWinCondition())

	// (+) Successful update of composite solution with two slots
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.updateGameState("17", 8)
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p1\n\n[_][2][3][_][5][6][_][8][9]", output)
	testing_utils.AssertEQb(t, false, stb.checkWinCondition())

	// (+) Successful update of composite solution with three slots
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.updateGameState("235", 10)
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p1\n\n[_][_][_][_][_][6][_][8][9]", output)
	testing_utils.AssertEQb(t, false, stb.checkWinCondition())

	// (+) Successful update of the last slot
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.updateGameState("9", 9)
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p1\n\n[_][_][_][_][_][6][_][8][_]", output)
	testing_utils.AssertEQb(t, false, stb.checkWinCondition())

	// (+) Successful update of internal isolated slot
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.updateGameState("8", 8)
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p1\n\n[_][_][_][_][_][6][_][_][_]", output)
	testing_utils.AssertEQb(t, false, stb.checkWinCondition())

	// (+) Successful update of the final open slot to close the box and win
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.updateGameState("6", 6)
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p1\n\n[_][_][_][_][_][_][_][_][_]", output)
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

func TestNextTurn(t *testing.T) {
	// Test the behavior for setting up the next turn, which includes:
	// - opening the box
	// - selecting the next player

	// Capture the print output for testing
	origStdout, r, w := testing_utils.RedirectStdout()
	stb := NewShutBox([]string{"p1", "p2", "p3", "p4"})
	stb.printGameState()
	output := testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p1\n\n[1][2][3][4][5][6][7][8][9]", output)

	// nextPlayer tests

	// Increment the player p1 -> p2
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.nextPlayer()
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p2\n\n[1][2][3][4][5][6][7][8][9]", output)

	// Increment the player p2 -> p3
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.nextPlayer()
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p3\n\n[1][2][3][4][5][6][7][8][9]", output)

	// Increment the player p3 -> p4
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.nextPlayer()
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p4\n\n[1][2][3][4][5][6][7][8][9]", output)

	// Increment the player p4 -> p1 LOOP BACK
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.nextPlayer()
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p1\n\n[1][2][3][4][5][6][7][8][9]", output)

	// resetBox tests

	// Close some slots
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.updateGameState("147", 12)
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p1\n\n[_][2][3][_][5][6][_][8][9]", output)

	// Increment the player p1 -> p2
	// Reset the box
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.nextPlayer()
	stb.resetBox()
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p2\n\n[1][2][3][4][5][6][7][8][9]", output)

	// Close some slots
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.updateGameState("147", 12)
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p2\n\n[_][2][3][_][5][6][_][8][9]", output)

	// Increment the player p1 -> p2
	// Reset the box
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.nextTurn()
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p3\n\n[1][2][3][4][5][6][7][8][9]", output)

	// Change player to p4 to test loop around
	stb.nextPlayer()
	// Close some slots
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.updateGameState("26", 8)
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p4\n\n[1][_][3][4][5][_][7][8][9]", output)

	// Increment the player p4 -> p1 LOOP BACK
	// Reset the box
	origStdout, r, w = testing_utils.RedirectStdout()
	stb.nextTurn()
	stb.printGameState()
	output = testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	testing_utils.AssertEQ(t, "Player: p1\n\n[1][2][3][4][5][6][7][8][9]", output)
}
