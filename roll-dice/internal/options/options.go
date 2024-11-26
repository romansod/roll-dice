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
)

/// Constants

const ErrUnsupported = "unsupported option"
const ErrNotImplemented = "not yet implemented"

/// Option Types

const (
	exit       = iota
	flip_coins = iota
	roll_dice  = iota
)

/// Collection of Options

type Options struct {
	opts map[int]Opt
}

func (options Options) displayOptions() {
	fmt.Print("\n\nPlease enter the option number\n\nRegistered Options:\n\n")
	for i := 0; i < len(options.opts); i++ {
		v := options.opts[i]
		fmt.Printf("\t%d) %s\n", v.getOptNum(), v.getName())
	}
}

/**
 * Register all Options
 *   opts[optNum] = Opt{name: _, optNum: _}
 */
func (options *Options) registerOptions() {
	options.opts = make(map[int]Opt)
	options.opts[exit] = OptExit{name: "Exit", optNum: exit}
	options.opts[flip_coins] = OptFlipCoins{name: "Flip Coins", optNum: flip_coins}
	options.opts[roll_dice] = OptRollDice{name: "Roll Dice", optNum: roll_dice}
}

/**
 * Process user input
 * params :
 *   stdin - holds user input
 * returns:
 *   menu  - option as number
 *   error - any error encountered by string to int conversion
 */
func (options Options) processInput(stdin io.Reader) (int, error) {
	options.displayOptions()
	scanner := bufio.NewScanner(stdin)
	scanner.Scan()
	fmt.Print("\n")
	return strconv.Atoi(scanner.Text())
}

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
	process() error
	getName() string
	getOptNum() int
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
	return errors.New(ErrNotImplemented)
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
	return errors.New(ErrNotImplemented)
}

func (optRollDice OptRollDice) getName() string {
	return optRollDice.name
}

func (optRollDice OptRollDice) getOptNum() int {
	return optRollDice.optNum
}

func Menu() {
	user_exit := false

	menu_options := Options{}
	menu_options.registerOptions()

	for !user_exit {
		// User input for menu option parsed
		input, err := menu_options.processInput(os.Stdin)
		if err != nil {
			fmt.Print("Syntax Error: Expected Integer\n")
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
