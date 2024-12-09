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
//	opts[optNum] = Opt{name: _, optNum: _}
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
//		error : any error encountered
func (options Options) runOption(opt int) error {
	opt_t, exists := options.opts[opt]
	if exists {
		return opt_t.process()
	} else {
		return errors.New(ErrUnsupported)
	}
}

/// - Base Opt type

type Opt interface {
	process() error  // Setup and execute operation
	getName() string // Retrieve the name of the operation
	getOptNum() int  // Get the opt number
}

/// - 0) Exit

type OptExit struct {
	name   string
	optNum int
}

func (optExit OptExit) process() error {
	fmt.Print("Exiting now ")
	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		fmt.Print(".")
	}

	fmt.Print("\n")
	return nil
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

func (optFlipCoins OptFlipCoins) process() error {
	// Only need the number of coin flips, the possible
	// outcomes are already known
	fmt.Print("Please enter the number of coin flips\n")
	input, err := processInput(os.Stdin)
	if err != nil {
		return errors.New(SyntaxErrExpectedInt)
	}

	coinFlip := probgen.CoinFlip{NumEvents: input}

	return probgen.ValidateAndExecute(coinFlip)
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

func (optRollDice OptRollDice) process() error {
	// Need the number of sides on the dice
	fmt.Printf("Please select the number of dice sides %s:\n", probgen.ValidDiceTypes)
	sides, err := processInput(os.Stdin)
	if err != nil {
		return errors.New(SyntaxErrExpectedInt)
	}

	// Need the number of rolls for the dice
	fmt.Print("Please enter the number of dice rolls\n")
	rolls, err := processInput(os.Stdin)
	if err != nil {
		return errors.New(SyntaxErrExpectedInt)
	}

	diceRoll := probgen.DiceRoll{NumEvents: rolls, NumSides: sides}

	return probgen.ValidateAndExecute(diceRoll)
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
//		int   : option as number
//		error : any error encountered by string to int conversion
func processInput(stdin io.Reader) (int, error) {
	scanner := bufio.NewScanner(stdin)
	scanner.Scan()
	// Add extra space after input to avoid clutter
	fmt.Print("\n")
	return strconv.Atoi(scanner.Text())
}

// Main driving function. Will continue to prompt user for input
// until failure or user asks to exit
func Menu() {
	user_exit := false

	menu_options := Options{}
	menu_options.registerOptions()

	for !user_exit {
		// User input for menu option parsed
		menu_options.displayOptions()
		input, err := processInput(os.Stdin)
		if err != nil {
			fmt.Print(SyntaxErrExpectedInt)
			continue
		}

		// For this iteration, we will run the selected option
		err = menu_options.runOption(input)
		if err != nil {
			fmt.Print(err)
		}

		// End of iteration
		// When exit is selected, stop here
		user_exit = input == exit
	}
}
