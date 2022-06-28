package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/liampulles/span-digital-ranking-cli/cmd/sportrank/internal/adapter"
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

type EngineImpl struct {
	rowIOGateway adapter.RowIOGateway
}

var _ Engine = &EngineImpl{}

func NewEngineImpl(rowIOGateway adapter.RowIOGateway) *EngineImpl {
	return &EngineImpl{
		rowIOGateway: rowIOGateway,
	}
}

func (ei *EngineImpl) Run(args []string, stdin io.Reader, stdout io.Writer) int {
	// Convert args to options.
	opts, err := ei.evaluateArgs(args, stdin, stdout)
	if err != nil {
		if errors.Is(err, errArgParse) {
			return FlagParseErrorCode
		}

		return ei.fail(err)
	}

	// Make sure we close anything that needs it.
	if closable, ok := opts.Input.(io.Closer); ok {
		defer closable.Close()
	}
	if closable, ok := opts.Output.(io.Closer); ok {
		defer closable.Close()
	}

	// Read input
	inputRows, err := ei.readLines(opts.Input)
	if err != nil {
		return ei.fail(err)
	}

	// Execute the business logic
	outputRows, err := ei.rowIOGateway.CalculateRankings(inputRows)
	if err != nil {
		return ei.fail(err)
	}

	// Write output
	if err := ei.writeLines(opts.Output, outputRows); err != nil {
		return ei.fail(err)
	}

	return SuccessCode
}

func (ei *EngineImpl) readLines(input io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		// Only include if not blank
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("could not parse input: %w", err)
	}
	return lines, nil
}

func (ei *EngineImpl) writeLines(output io.Writer, lines []string) error {
	// Use a buffered writer for predictable performance
	bufOutput := bufio.NewWriter(output)
	for _, line := range lines {
		if _, err := fmt.Fprintln(bufOutput, line); err != nil {
			return fmt.Errorf("could not write to output: %w", err)
		}
	}
	if err := bufOutput.Flush(); err != nil {
		return fmt.Errorf("could not flush output: %w", err)
	}
	return nil
}

func (ei *EngineImpl) fail(err error) int {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	return ei.chooseExitCode(err)
}

func (ei *EngineImpl) chooseExitCode(err error) int {
	if errors.Is(err, errCouldNotOpenInput) {
		return CouldNotReadInputCode
	}
	if errors.Is(err, errCouldNotOpenOutput) {
		return CouldNotWriteOutputCode
	}
	if errors.Is(err, adapter.ErrMalformedRow) {
		return InvalidFormatCode
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
