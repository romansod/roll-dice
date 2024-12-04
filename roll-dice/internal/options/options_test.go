package options

import (
	"bytes"
	"testing"

	"github.com/romansod/roll-dice/internal/testing_utils"
)

// Initialize Options and registers them for use during tests
//
//	Returns
//		options : the initialized and registered Options
func setUp() Options {
	options := Options{}
	options.registerOptions()

	return options
}

func TestDisplay(t *testing.T) {
	// Tests the formatting of the commandline display
	// This will likely need updates for every new feature

	options := setUp()
	origStdout, r, w := testing_utils.RedirectStdout()

	expected := ""
	options.displayOptions()

	output := testing_utils.CaptureOutput(r, w, origStdout)
	expected =
		"\n\nPlease enter the option number" +
			"\n\nRegistered Options:\n" +
			"\n\t0) Exit" +
			"\n\t1) Flip Coins" +
			"\n\t2) Roll Dice\n"
	if !testing_utils.AssertEQ(expected, output) {
		t.Errorf(testing_utils.AssertFailed, expected, output)
	}
}

func TestProcessInput(t *testing.T) {
	// Tests input processing for error and passing values

	options := setUp()
	origStdout, ignoreOut := testing_utils.IgnoreStdout()

	expected := ""
	var stdin bytes.Buffer

	// Err : non-numeric
	stdin.Write([]byte("invalid"))
	expected = "strconv.Atoi: parsing \"invalid\": invalid syntax"
	_, err := options.processInput(&stdin)
	if !testing_utils.AssertEQ(expected, err.Error()) {
		t.Errorf(testing_utils.AssertFailed, expected, err)
	}
	stdin.Reset()

	// Pass : unsupported
	stdin.Write([]byte("-1"))
	input, err := options.processInput(&stdin)
	if !testing_utils.AssertNIL(err) {
		t.Errorf(testing_utils.AssertFailed, "nil", err)
	} else if !testing_utils.AssertEQi(-1, input) {
		t.Errorf(testing_utils.AssertFailed, -1, input)
	}
	stdin.Reset()

	testing_utils.IgnoreStdoutClose(origStdout, ignoreOut)
}

func TestMenuOptions(t *testing.T) {
	// Tests the selection of menu options
	// This will likely need updates for every new feature

	options := setUp()
	origStdout, ignoreOut := testing_utils.IgnoreStdout()

	expected := ""

	/// - -1) Unsupported Option
	err := options.runOption(-1)
	expected = "unsupported option"
	if !testing_utils.AssertEQ(expected, err.Error()) {
		t.Errorf(testing_utils.AssertFailed, expected, err)
	}

	/// - 0) Exit
	err = options.runOption(exit)
	expected = "nil"
	if !testing_utils.AssertNIL(err) {
		t.Errorf(testing_utils.AssertFailed, expected, err)
	}

	/// - 1) Flip Coins
	err = options.runOption(flip_coins)
	expected = "not yet implemented"
	if !testing_utils.AssertEQ(expected, err.Error()) {
		t.Errorf(testing_utils.AssertFailed, expected, err)
	}

	/// - 2) Roll Dice
	err = options.runOption(roll_dice)
	expected = "not yet implemented"
	if !testing_utils.AssertEQ(expected, err.Error()) {
		t.Errorf(testing_utils.AssertFailed, expected, err)
	}

	testing_utils.IgnoreStdoutClose(origStdout, ignoreOut)
}
