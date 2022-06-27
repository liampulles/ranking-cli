package wire_test

import (
	"testing"

	"github.com/liampulles/span-digital-ranking-cli/cmd/sportrank/internal/wire"
	"github.com/stretchr/testify/assert"
)

func ExampleRun_validCase() {
	wire.Run([]string{"prog.name", "-i", "testdata/valid_input.txt"})
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

	// Exercise SUT
	actualCode := wire.Run(argsFixture)

	// Verify results
	assert.Equal(t, wire.SuccessCode, actualCode)
}

func TestRun_GivenInvalidInput_ShouldReturnInvalidFormat(t *testing.T) {
	// Setup fixture
	argsFixture := []string{"prog.name", "-i", "testdata/invalid_input.txt"}

	// Exercise SUT
	actualCode := wire.Run(argsFixture)

	// Verify results
	assert.Equal(t, wire.InvalidFormat, actualCode)
}

func TestRun_GivenInputDoesNotExist_ShouldReturnCouldNotReadInput(t *testing.T) {
	// Setup fixture
	argsFixture := []string{"prog.name", "-i", "does.not.exist.txt"}

	// Exercise SUT
	actualCode := wire.Run(argsFixture)

	// Verify results
	assert.Equal(t, wire.CouldNotReadInput, actualCode)
}
