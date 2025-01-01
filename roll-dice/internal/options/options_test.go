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

	output := testing_utils.CaptureAndRestoreOutput(r, w, origStdout)
	expected =
		"\n\nPlease enter the option number" +
			"\n\nRegistered Options:\n" +
			"\n\t0) Exit" +
			"\n\t1) Flip Coins" +
			"\n\t2) Roll Dice\n"
	testing_utils.AssertEQ(t, expected, output)
}

func TestProcessInput(t *testing.T) {
	// Tests input processing for error and passing values
	// Ignoring stdout helps with extra lines added to processInput
	origStdout, ignoreOut := testing_utils.IgnoreStdout()

	expected := ""
	var stdin bytes.Buffer

	// Err : non-numeric
	stdin.Write([]byte("invalid"))
	expected = "strconv.Atoi: parsing \"invalid\": invalid syntax"
	_, _, err := processInput(&stdin)
	testing_utils.AssertEQ(t, expected, err.Error())
	stdin.Reset()

	// Pass : unsupported
	stdin.Write([]byte("-1"))
	_, input, err := processInput(&stdin)
	testing_utils.AssertNIL(t, err)
	testing_utils.AssertEQi(t, -1, input)
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
	done, err := options.runOption(-1)
	expected = ErrUnsupported
	testing_utils.AssertEQ(t, expected, err.Error())
	testing_utils.AssertEQb(t, false, done)

	/// - 0) Exit
	done, err = options.runOption(exit)
	expected = "nil"
	testing_utils.AssertNIL(t, err)
	testing_utils.AssertEQb(t, true, done)

	/// - 1) Flip Coins
	done, err = options.runOption(flip_coins)
	testing_utils.AssertNIL(t, err)
	testing_utils.AssertEQb(t, false, done)

	/// - 2) Roll Dice
	done, err = options.runOption(roll_dice)
	testing_utils.AssertNIL(t, err)
	testing_utils.AssertEQb(t, false, done)

	testing_utils.IgnoreStdoutClose(origStdout, ignoreOut)
}
