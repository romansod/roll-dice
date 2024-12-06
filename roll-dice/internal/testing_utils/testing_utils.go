package testing_utils

/*
	To run unit tests, run from the roll-dice directory:

		go test ./... -v
*/
import (
	"bytes"
	"os"
	"strconv"
	"testing"
)

const AssertFailed = "\nexpected : %v\nactual   : %v\n"

// Disable stdout for test purposes. Called in conjunction
// with RestoreStdout
//
//	Returns
//		orig     : original stdout to restore later
//		nullFile : temporary null file to close
func IgnoreStdout() (*os.File, *os.File) {
	origStdout := os.Stdout
	nullFile, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
	os.Stdout = nullFile
	return origStdout, nullFile
}

// Restore stdout for test purposes and close nullFile.
// Called in conjunction with IgnoreStdout
//
//	Params
//		orig     : original stdout to restore
//		nullFile : temporary null file to close
func IgnoreStdoutClose(orig *os.File, nullFile *os.File) {
	os.Stdout = orig
	nullFile.Close()
}

// Create and redirect stdout to a file variable with pipe
// and return original stdout destination so it can be reset
// later with CaptureAndRestoreOutput
//
//	Returns
//		*os.File : original stdout for restoring later
//		*os.File : read end of pipe
//		*os.File : write end of pipe
func RedirectStdout() (*os.File, *os.File, *os.File) {
	// Save the original os.Stdout
	origStdout := os.Stdout

	// Create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	// Redirect os.Stdout to the pipe writer
	os.Stdout = w

	return origStdout, r, w
}

// Capture output from pipe, then restore writer to previous destination
//
//	Params
//		r *os.File : Read end of pipe
//		w *os.File : Write end of pipe
//		origStdout *os.File : original write end of pipe to restore to
//	Returns
//		string : Contents of r -> w
func CaptureAndRestoreOutput(r *os.File, w *os.File, origStdout *os.File) string {
	// Close the writer to signal end of writing
	w.Close()

	// Read the output from the pipe
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		panic(err)
	}

	// Restore the original os.Stdout
	os.Stdout = origStdout

	return buf.String()
}

// Assert that Expected == Actual. If false then
// report an error
//
//	Params
//		t *testing.T : needed for calling Errorf
//		exp string   : Expected value
//		act string   : Actual value
func AssertEQ(t *testing.T, exp string, act string) {
	if exp != act {
		t.Errorf(AssertFailed, exp, act)
	}
}

// Assert that Expected == Actual. If false then
// report an error
//
//	Params
//		t *testing.T : needed for calling Errorf
//		exp int      : Expected value
//		act int      : Actual value
func AssertEQi(t *testing.T, exp int, act int) {
	AssertEQ(t, strconv.Itoa(exp), strconv.Itoa(act))
}

// Assert that err is nil, ie no error occurred. If false then
// report an error
//
//	Params
//		t *testing.T : needed for calling Errorf
//		err error    : Error variable to check if nil
func AssertNIL(t *testing.T, err error) {
	if err != nil {
		t.Errorf(AssertFailed, "nil", err.Error())
	}
}
