package wire

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Exit codes:
const (
	// Note: Don't use iota here, we don't want these to change dynamically over time
	// as external tools may depend on these values.
	SuccessCode         = 0
	InvalidFormat       = 1
	CouldNotReadInput   = 2
	CouldNotWriteOutput = 3
	InternalError       = 4
	FlagParseError      = 5
)

func Run(args []string, stdin io.Reader, stdout io.Writer) int {
	// Convert args to options.
	opts, err := evaluateArgs(args, stdin, stdout)
	if err != nil {
		if errors.Is(err, ErrArgParse) {
			return FlagParseError
		}

		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return chooseExitCode(err)
	}

	// Make sure we close anything that needs it.
	if closable, ok := opts.Input.(io.Closer); ok {
		defer closable.Close()
	}
	if closable, ok := opts.Output.(io.Closer); ok {
		defer closable.Close()
	}

	fmt.Fprintf(opts.Output, "%+v\n", opts)

	return SuccessCode
}

func chooseExitCode(err error) int {
	if errors.Is(err, ErrCouldNotOpenInput) {
		return CouldNotReadInput
	}
	if errors.Is(err, ErrCouldNotOpenOutput) {
		return CouldNotWriteOutput
	}
	return InternalError
}

// --- Argument related ---

// Defined argument related errors:
var (
	ErrArgParse           = errors.New("arg parsing error")
	ErrCouldNotOpenInput  = errors.New("could not open input")
	ErrCouldNotOpenOutput = errors.New("could not open output")
)

type options struct {
	Input  io.Reader
	Output io.Writer
}

func evaluateArgs(args []string, stdin io.Reader, stdout io.Writer) (options, error) {
	// Define and run the flag set.
	flagSet := flag.NewFlagSet("sportrank", flag.ContinueOnError)
	inputPtr := flagSet.String("i", "-", "Input file, or - for STDIN.")
	outputPtr := flagSet.String("o", "-", "Output file, or - for STDOUT.")
	if err := flagSet.Parse(args[1:]); err != nil {
		return options{}, ErrArgParse
	}

	// Get the appropriate input and output, given the args.
	input, err := getFileSource(*inputPtr, stdin, os.O_RDONLY, 0777, ErrCouldNotOpenInput)
	if err != nil {
		flagSet.Usage()
		return options{}, err
	}
	output, err := getFileSource(*outputPtr, stdout, os.O_RDWR|os.O_CREATE, 0755, ErrCouldNotOpenOutput)
	if err != nil {
		flagSet.Usage()
		return options{}, err
	}

	return options{
		Input:  input.(io.Reader),
		Output: output.(io.Writer),
	}, nil
}

func getFileSource(
	arg string,
	std interface{},
	flag int,
	createPerm os.FileMode,
	couldNotOpenErr error,
) (interface{}, error) {
	cleaned := strings.TrimSpace(arg)

	if cleaned == "-" {
		return std, nil
	}

	f, err := os.OpenFile(cleaned, flag, createPerm)
	if err != nil {
		return nil, fmt.Errorf("%s - %w", cleaned, couldNotOpenErr)
	}
	return f, nil
}
