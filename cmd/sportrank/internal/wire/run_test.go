package wire_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/liampulles/span-digital-ranking-cli/cmd/sportrank/internal/wire"
	"github.com/stretchr/testify/assert"
)

func ExampleRun_validCase() {
	wire.Run([]string{"prog.name", "-i", "testdata/valid_input.txt"}, os.Stdin, os.Stdout)
	// Output:
	// 1. Tarantulas, 6 pts
	// 2. Lions, 5 pts
	// 3. FC Awesome, 1 pt
	// 3. Snakes, 1 pt
	// 5. Grouches, 0 pts
}

func TestRun_GivenValidInput_ShouldReturnSuccess(t *testing.T) {
	// Setup fixture
	argsFixture := []string{"prog.name", "-i", "testdata/valid_input.txt"}
	output := bytes.NewBufferString("")

	// Setup expectations
	expectedOutput := `1. Tarantulas, 6 pts
2. Lions, 5 pts
3. FC Awesome, 1 pt
3. Snakes, 1 pt
5. Grouches, 0 pts`

	// Exercise SUT
	actualCode := wire.Run(argsFixture, nil, output)

	// Verify results
	assert.Equal(t, wire.SuccessCode, actualCode)
	assert.Equal(t, expectedOutput, output.String())
}

func TestRun_GivenValidInputViaStdin_ShouldReturnSuccess(t *testing.T) {
	// Setup fixture
	argsFixture := []string{"prog.name", "-i", "-"}
	input, err := os.Open("testdata/valid_input.txt")
	if err != nil {
		t.FailNow()
	}
	output := bytes.NewBufferString("")

	// Setup expectations
	expectedOutput := `1. Tarantulas, 6 pts
2. Lions, 5 pts
3. FC Awesome, 1 pt
3. Snakes, 1 pt
5. Grouches, 0 pts`

	// Exercise SUT
	actualCode := wire.Run(argsFixture, input, output)

	// Verify results
	assert.Equal(t, wire.SuccessCode, actualCode)
	assert.Equal(t, expectedOutput, output.String())
}

func TestRun_GivenValidInputAndOutputViaFile_ShouldReturnSuccess(t *testing.T) {
	// Setup fixture
	argsFixture := []string{"prog.name", "-i", "testdata/valid_input.txt", "-o", "testdata/temptestout.txt"}
	// Setup expectations
	expectedOutput := `1. Tarantulas, 6 pts
2. Lions, 5 pts
3. FC Awesome, 1 pt
3. Snakes, 1 pt
5. Grouches, 0 pts`

	// Exercise SUT
	actualCode := wire.Run(argsFixture, nil, os.Stdout)

	// Verify results
	assert.Equal(t, wire.SuccessCode, actualCode)
	outputBytes, err := ioutil.ReadFile("testdata/temptestout.txt")
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, expectedOutput, string(outputBytes))
}

func TestRun_GivenInvalidInput_ShouldReturnInvalidFormat(t *testing.T) {
	// Setup fixture
	argsFixture := []string{"prog.name", "-i", "testdata/invalid_input.txt"}

	// Exercise SUT
	actualCode := wire.Run(argsFixture, nil, os.Stdout)

	// Verify results
	assert.Equal(t, wire.InvalidFormat, actualCode)
}

func TestRun_GivenInputDoesNotExist_ShouldReturnCouldNotReadInput(t *testing.T) {
	// Setup fixture
	argsFixture := []string{"prog.name", "-i", "does.not.exist.txt"}

	// Exercise SUT
	actualCode := wire.Run(argsFixture, nil, os.Stdout)

	// Verify results
	assert.Equal(t, wire.CouldNotReadInput, actualCode)
}

func TestRun_GivenCouldNotOpenOutput_ShouldReturnCouldNotWriteOutput(t *testing.T) {
	// Setup fixture
	// -> We make use of the fact that you cannot
	//    create a file with a name longer than 255 chars.
	argsFixture := []string{"prog.name", "-o", strings.Repeat("x", 2000)}

	// Exercise SUT
	actualCode := wire.Run(argsFixture, os.Stdin, nil)

	// Verify results
	assert.Equal(t, wire.CouldNotWriteOutput, actualCode)
}

func TestRun_GivenInvalidArgs_ShouldReturnFlagParseError(t *testing.T) {
	// Setup fixture
	argsFixture := []string{"prog.name", "-b"}

	// Exercise SUT
	actualCode := wire.Run(argsFixture, os.Stdin, nil)

	// Verify results
	assert.Equal(t, wire.FlagParseError, actualCode)
}
