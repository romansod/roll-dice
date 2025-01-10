package main

import (
	"fmt"
	"os"

	"github.com/romansod/roll-dice/internal/options"
)

const instructions string = "\nSelect the menu option using the associated\n" +
	"integer. Additionally, an empty input\n" +
	"indicates you are 'done' while executing an\n" +
	"operation, returning execution to the main menu\n\n"

func main() {
	fmt.Print("--------------- Welcome ---------------\n")
	fmt.Print(instructions)

	options.Menu()
	os.Exit(0)
}
