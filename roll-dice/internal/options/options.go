/*
options.go - utility

Called by driver.go through Menu()

Handles all input validation for the menu
*/
package options

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/romansod/roll-dice/internal/probgen"
)

/// Constants

const ErrUnsupported = "unsupported option"
const ErrNotImplemented = "not yet implemented"

const SyntaxErrExpectedInt = "syntax error: expected integer"

/// Option Types

const (
	exit       = iota
	flip_coins = iota
	roll_dice  = iota
)

/// Collection of Options

type Options struct {
	opts map[int]Opt // Map of menu options to Opt
}

// Print the menu options
func (options Options) displayOptions() {
	fmt.Print("\n\nPlease enter the option number\n\nRegistered Options:\n\n")
	for i := 0; i < len(options.opts); i++ {
		v := options.opts[i]
		fmt.Printf("\t%d) %s\n", v.getOptNum(), v.getName())
	}
}

// Register all Options
//
//	Example:
//		opts[optNum] = Opt{name: _, optNum: _}
func (options *Options) registerOptions() {
	options.opts = make(map[int]Opt)
	options.opts[exit] = OptExit{name: "Exit", optNum: exit}
	options.opts[flip_coins] = OptFlipCoins{name: "Flip Coins", optNum: flip_coins}
	options.opts[roll_dice] = OptRollDice{name: "Roll Dice", optNum: roll_dice}
}

// Run the given Opt based on the opt number provided
//
//	Params
//		opt int : the menu option dictating which Opt is run
//	Returns
//		bool  : true if user indicates they are done
//		error : any error encountered
func (options Options) runOption(opt int) (bool, error) {
	done, err := false, errors.New(ErrUnsupported)
	opt_t, exists := options.opts[opt]
	if exists {
		for !done {
			done, err = opt_t.process()
		}

		done = done && opt == exit
	}

	return done, err
}

/// - Base Opt type

type Opt interface {
	process() (bool, error) // Setup and execute operation
	getName() string        // Retrieve the name of the operation
	getOptNum() int         // Get the opt number
}

/// - 0) Exit

type OptExit struct {
	name   string
	optNum int
}

func (optExit OptExit) process() (bool, error) {
	// Simply exit gracefully

	fmt.Print("Exiting now ")
	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		fmt.Print(".")
	}

	fmt.Print("\n")
	return true, nil
}

func (optExit OptExit) getName() string {
	return optExit.name
}

func (optExit OptExit) getOptNum() int {
	return optExit.optNum
}

/// - 1) Flip Coins

type OptFlipCoins struct {
	name   string
	optNum int
}

func (optFlipCoins OptFlipCoins) process() (bool, error) {
	// Prompt user for the number of coin flips they want to do

	fmt.Print("Please enter the number of coin flips:\n")
	done, input, err := processInput(os.Stdin)

	if done {
		return true, err
	}

	if err != nil {
		return false, errors.New(SyntaxErrExpectedInt)
	}

	coinFlip := probgen.NewCoinFlip(input)

	return false, probgen.ValidateAndExecute(coinFlip)
}

func (optFlipCoins OptFlipCoins) getName() string {
	return optFlipCoins.name
}

func (optFlipCoins OptFlipCoins) getOptNum() int {
	return optFlipCoins.optNum
}

/// - 2) Roll Dice

type OptRollDice struct {
	name   string
	optNum int
}

func (optRollDice OptRollDice) process() (bool, error) {
	// Prompt the user for the number of sides on the dice
	fmt.Printf("Please select the number of dice sides %s:\n", probgen.ValidDiceTypes)
	done, sides, err := processInput(os.Stdin)
	if done {
		return true, err
	}

	if err != nil {
		return false, errors.New(SyntaxErrExpectedInt)
	}

	// Prompt the user for the number of rolls for the dice
	fmt.Print("Please enter the number of dice rolls:\n")
	done, rolls, err := processInput(os.Stdin)
	if done {
		return true, err
	}

	if err != nil {
		return false, errors.New(SyntaxErrExpectedInt)
	}

	diceRoll := probgen.NewDiceRoll(rolls, sides)

	return false, probgen.ValidateAndExecute(diceRoll)
}

func (optRollDice OptRollDice) getName() string {
	return optRollDice.name
}

func (optRollDice OptRollDice) getOptNum() int {
	return optRollDice.optNum
}

// Process user input
//
//	Params
//		stdin io.Reader : holds user input
//
//	Returns
//		bool  : true if user indicates they are done
//		int   : option as number
//		error : any error encountered by string to int conversion
func processInput(stdin io.Reader) (bool, int, error) {
	scanner := bufio.NewScanner(stdin)
	scanner.Scan()
	if scanner.Text() == "" {
		// User is done providing inputs
		return true, -1, nil
	}
	// Add extra space after input to avoid clutter
	fmt.Print("\n")
	input, err := strconv.Atoi(scanner.Text())
	return false, input, err
}

// Main driving function. Will continue to prompt user for input
// until failure or user asks to exit
func Menu() {
	done, input, err := false, -1, error(nil)

	menu_options := Options{}
	menu_options.registerOptions()

	// Consumer user input until user is done and indicates exit
	// through the Exit option
	for !done {
		// User input for menu option parsed
		menu_options.displayOptions()
		// Errors from processing options fall back to the
		// main menu to here where user is prompted again
		_, input, err = processInput(os.Stdin)

		if err != nil {
			fmt.Print(SyntaxErrExpectedInt)
			continue
		}

		// For this iteration, we will run the selected option
		done, err = menu_options.runOption(input)
		if err != nil {
			fmt.Print(err)
		}
	}
}
