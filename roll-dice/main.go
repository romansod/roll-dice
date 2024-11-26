package main

import (
	"fmt"
	"os"

	"github.com/romansod/roll-dice/internal/options"
)

func main() {
	fmt.Print("--------------- Welcome ---------------\n")

	options.Menu()
	os.Exit(0)
}
