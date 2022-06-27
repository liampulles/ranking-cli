package cli

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
	SuccessCode             = 0
	InvalidFormatCode       = 1
	CouldNotReadInputCode   = 2
	CouldNotWriteOutputCode = 3
	InternalErrorCode       = 4
	FlagParseErrorCode      = 5
)

// Engine facilitates control of the system via a CLI.
type Engine interface {
	Run(args []string, stdin io.Reader, stdout io.Writer) int
}

type EngineImpl struct{}

var _ Engine = &EngineImpl{}

func NewEngineImpl() *EngineImpl {
	return &EngineImpl{}
}

func (ei *EngineImpl) Run(args []string, stdin io.Reader, stdout io.Writer) int {
	// Convert args to options.
	opts, err := ei.evaluateArgs(args, stdin, stdout)
	if err != nil {
		if errors.Is(err, errArgParse) {
			return FlagParseErrorCode
		}

		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return ei.chooseExitCode(err)
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

func (ei *EngineImpl) chooseExitCode(err error) int {
	if errors.Is(err, errCouldNotOpenInput) {
		return CouldNotReadInputCode
	}
	if errors.Is(err, errCouldNotOpenOutput) {
		return CouldNotWriteOutputCode
	}
	return InternalErrorCode
}

// --- Argument related ---

// Defined argument related errors:
var (
	errArgParse           = errors.New("arg parsing error")
	errCouldNotOpenInput  = errors.New("could not open input")
	errCouldNotOpenOutput = errors.New("could not open output")
)

type options struct {
	Input  io.Reader
	Output io.Writer
}

func (ei *EngineImpl) evaluateArgs(args []string, stdin io.Reader, stdout io.Writer) (options, error) {
	// Define and run the flag set.
	flagSet := flag.NewFlagSet("sportrank", flag.ContinueOnError)
	inputPtr := flagSet.String("i", "-", "Input file, or - for STDIN.")
	outputPtr := flagSet.String("o", "-", "Output file, or - for STDOUT.")
	if err := flagSet.Parse(args[1:]); err != nil {
		return options{}, errArgParse
	}

	// Get the appropriate input and output, given the args.
	input, err := ei.getFileSource(*inputPtr, stdin, os.O_RDONLY, 0777, errCouldNotOpenInput)
	if err != nil {
		flagSet.Usage()
		return options{}, err
	}
	output, err := ei.getFileSource(*outputPtr, stdout, os.O_RDWR|os.O_CREATE, 0755, errCouldNotOpenOutput)
	if err != nil {
		flagSet.Usage()
		return options{}, err
	}

	return options{
		Input:  input.(io.Reader),
		Output: output.(io.Writer),
	}, nil
}

func (ei *EngineImpl) getFileSource(
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
