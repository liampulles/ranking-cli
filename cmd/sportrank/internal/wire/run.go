package wire

import (
	"fmt"
)

// Exit codes
const (
	// Note: Don't use iota here, we don't want these to change dynamically over time
	// as external tools may depend on these values.
	SuccessCode       = 0
	InvalidFormat     = 1
	CouldNotReadInput = 2
)

func Run(args []string) int {
	fmt.Println("Hello world!")
	return SuccessCode
}
