package testing_utils

import (
	"bytes"
	"os"
	"strconv"
)

const AssertFailed = "\nexpected : %v\nactual   : %v\n"

/**
 * Disable stdout for test purposes. Called in conjunction
 * with restoreStdout
 *
 * params :
 *   NONE
 * returns:
 *   original stdout - to restore later
 *   nullFile        - to close later
 */
func IgnoreStdout() (*os.File, *os.File) {
	origStdout := os.Stdout
	nullFile, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
	os.Stdout = nullFile
	return origStdout, nullFile
}

/**
 * Restore stdout for test purposes and close nullFile.
 * Called in conjunction with ignoreStdout
 *
 * params :
 *   orig     - original stdout to restore
 *   nullFile - temporary null file to close
 * returns:
 *   NONE
 */
func IgnoreStdoutClose(orig *os.File, nullFile *os.File) {
	os.Stdout = orig
	nullFile.Close()
}

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

func CaptureOutput(r *os.File, w *os.File, origStdout *os.File) string {
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

func AssertEQ(exp string, act string) bool {
	return exp == act
}

func AssertEQi(exp int, act int) bool {
	return AssertEQ(strconv.Itoa(exp), strconv.Itoa(act))
}

func AssertNIL(err error) bool {
	return err == nil
}
