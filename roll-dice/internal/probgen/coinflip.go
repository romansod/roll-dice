/*
coinflip.go

CoinFlip is a ProbEventType which
describes coin flips
*/
package probgen

import "fmt"

// Potential values
const (
	Heads = "Heads"
	Tails = "Tails"
)

type CoinFlip struct {
	numEvents int // number of coin flips
}

// Initialize private fields
//
//	Params
//		nEvents int : number of CoinFlip events
//	Returns
//		*CoinFlip : new CoinFlip object
func NewCoinFlip(nEvents int) *CoinFlip {
	return &CoinFlip{
		numEvents: nEvents,
	}
}

func (coinFlip CoinFlip) validate() (bool, error) {
	// Nothing to do as Heads and Tails are already implied
	return true, nil
}

func (coinFlip CoinFlip) execute() error {
	res, err := GenerateProbabilisticEvent(
		coinFlip.numEvents,
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
// numEvents: 123435
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
		Percent(res[Heads], coinFlip.numEvents), res[Heads],
		Percent(res[Tails], coinFlip.numEvents), res[Tails])
}

// Retrieve number of events
//
//	Returns
//		int : number of events
func (coinFlip CoinFlip) getNumEvents() int {
	return coinFlip.numEvents
}
