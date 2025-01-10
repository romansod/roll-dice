package utilities

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Process user number input
//
//	Params
//		stdin io.Reader : holds user input
//
//	Returns
//		bool  : true if user indicates they are done
//		int   : option as number
//		error : any error encountered by string to int conversion
func ProcessInputInt(stdin io.Reader) (bool, int, error) {
	done, input_str := ProcessInputStr(stdin)
	input_i := -1
	if done {
		return true, -1, nil
	}

	input_i, err := strconv.Atoi(input_str)

	return false, input_i, err
}

// Process user string input
//
//	Params
//		stdin io.Reader : holds user input
//
//	Returns
//		bool   : true if user indicates they are done
//		string : option as number
func ProcessInputStr(stdin io.Reader) (bool, string) {
	scanner := bufio.NewScanner(stdin)
	scanner.Scan()
	if scanner.Text() == "" {
		// User is done providing inputs
		fmt.Print("Stopping current operation\n")
		return true, ""
	}
	// Add extra space after input to avoid clutter
	fmt.Print("\n")
	return false, scanner.Text()
}
