/*
coinflip.go

CoinFlip is a ProbEventType which
describes coin flips
*/
package probgen

import "fmt"

const Heads = "Heads"
const Tails = "Tails"

type CoinFlip struct {
	NumEvents int // number of coin flips
}

func (coinFlip CoinFlip) validate() (bool, error) {
	// Nothing to do as Heads and Tails are already implied
	return true, nil
}

func (coinFlip CoinFlip) execute() error {
	res, err := GenerateProbabilisticEvent(
		coinFlip.NumEvents,
		[]string{
			Heads,
			Tails})

	if err == nil {
		coinFlip.display(res)
	}

	return err
}

// Print the coin flip results. Example:
//
// NumEvents: 123435
//
// (H) :  49.882126% : 61572
//
// (T) :  50.117874% : 61863
//
//	Params
//		res map[string]int : results of coin flips
func (coinFlip CoinFlip) display(res map[string]int) {
	fmt.Printf(
		"(H) : %10f%% : %d\n(T) : %10f%% : %d\n",
		Percent(res[Heads], coinFlip.NumEvents), res[Heads],
		Percent(res[Tails], coinFlip.NumEvents), res[Tails])
}
