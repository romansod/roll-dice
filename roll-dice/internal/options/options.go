/*
options.go - utility

Called by driver.go through Menu()

Handles all input validation for the menu
*/
package options

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/romansod/roll-dice/internal/games"
	"github.com/romansod/roll-dice/internal/probgen"
	"github.com/romansod/roll-dice/internal/utilities"
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
	shutthebox = iota
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
	options.opts[shutthebox] = OptShutTheBox{name: "Shut the Box", optNum: shutthebox}
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

			if err != nil {
				// Give feedback on any errors before next prompt
				fmt.Print(err.Error())
			}

			fmt.Print("\n\n")
		}

		fmt.Print("Returning to main menu ...\n")

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
	done, input, err := utilities.ProcessInputInt(os.Stdin)

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
	done, sides, err := utilities.ProcessInputInt(os.Stdin)
	if done {
		return true, err
	}

	if err != nil {
		return false, errors.New(SyntaxErrExpectedInt)
	}

	// Prompt the user for the number of rolls for the dice
	fmt.Print("Please enter the number of dice rolls:\n")
	done, rolls, err := utilities.ProcessInputInt(os.Stdin)
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

/// - 3) Roll Dice

type OptShutTheBox struct {
	name   string
	optNum int
}

func (optShutTheBox OptShutTheBox) process() (bool, error) {
	done, players, err := getPlayers(os.Stdin)
	if done {
		return true, err
	}

	if err != nil {
		return false, err
	}

	shutTheBox := games.NewShutBox(players)
	shutTheBox.Run()

	return true, nil
}

func (optShutTheBox OptShutTheBox) getName() string {
	return optShutTheBox.name
}

func (optShutTheBox OptShutTheBox) getOptNum() int {
	return optShutTheBox.optNum
}

func getPlayers(stdin io.Reader) (bool, []string, error) {
	// Prompt the user for the number of sides on the dice
	fmt.Print("Please indicate the number of players:\n")
	done, players_n, err := utilities.ProcessInputInt(stdin)
	if done {
		return true, nil, err
	}

	if err != nil {
		return false, nil, errors.New(SyntaxErrExpectedInt)
	}

	players := make([]string, players_n)

	for i := 0; i < players_n; i++ {
		// Prompt the user for the number of rolls for the dice
		fmt.Printf("Please enter player %d's name:\n", i+1)
		done, player := utilities.ProcessInputStr(stdin)
		if done {
			return true, nil, nil
		}

		players[i] = player
	}

	return false, players, nil
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
		_, input, err = utilities.ProcessInputInt(os.Stdin)

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
