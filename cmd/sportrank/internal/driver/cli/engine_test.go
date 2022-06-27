package cli_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/liampulles/span-digital-ranking-cli/cmd/sportrank/internal/driver/cli"
	"github.com/stretchr/testify/suite"
)

type EngineImplTestSuite struct {
	suite.Suite
	sut *cli.EngineImpl
}

func TestEngineImplTestSuite(t *testing.T) {
	suite.Run(t, new(EngineImplTestSuite))
}

func (suite *EngineImplTestSuite) SetupTest() {
	suite.sut = cli.NewEngineImpl()
}

func (suite *EngineImplTestSuite) TestRun_GivenValidInput_ShouldReturnSuccess() {
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
	actualCode := suite.sut.Run(argsFixture, nil, output)

	// Verify results
	suite.Equal(cli.SuccessCode, actualCode)
	suite.Equal(expectedOutput, output.String())
}

func (suite *EngineImplTestSuite) TestRun_GivenValidInputViaStdin_ShouldReturnSuccess() {
	// Setup fixture
	argsFixture := []string{"prog.name", "-i", "-"}
	input, err := os.Open("testdata/valid_input.txt")
	if err != nil {
		suite.FailNow("could not read test input, please check")
	}
	output := bytes.NewBufferString("")

	// Setup expectations
	expectedOutput := `1. Tarantulas, 6 pts
2. Lions, 5 pts
3. FC Awesome, 1 pt
3. Snakes, 1 pt
5. Grouches, 0 pts`

	// Exercise SUT
	actualCode := suite.sut.Run(argsFixture, input, output)

	// Verify results
	suite.Equal(cli.SuccessCode, actualCode)
	suite.Equal(expectedOutput, output.String())
}

func (suite *EngineImplTestSuite) TestRun_GivenValidInputAndOutputViaFile_ShouldReturnSuccess() {
	// Setup fixture
	argsFixture := []string{"prog.name", "-i", "testdata/valid_input.txt", "-o", "testdata/temptestout.txt"}
	// Setup expectations
	expectedOutput := `1. Tarantulas, 6 pts
2. Lions, 5 pts
3. FC Awesome, 1 pt
3. Snakes, 1 pt
5. Grouches, 0 pts`

	// Exercise SUT
	actualCode := suite.sut.Run(argsFixture, nil, os.Stdout)

	// Verify results
	suite.Equal(cli.SuccessCode, actualCode)
	outputBytes, err := ioutil.ReadFile("testdata/temptestout.txt")
	suite.NoError(err)
	suite.Equal(expectedOutput, string(outputBytes))
}

func (suite *EngineImplTestSuite) TestRun_GivenInvalidInput_ShouldReturnInvalidFormat() {
	// Setup fixture
	argsFixture := []string{"prog.name", "-i", "testdata/invalid_input.txt"}

	// Exercise SUT
	actualCode := suite.sut.Run(argsFixture, nil, os.Stdout)

	// Verify results
	suite.Equal(cli.InvalidFormatCode, actualCode)
}

func (suite *EngineImplTestSuite) TestRun_GivenInputDoesNotExist_ShouldReturnCouldNotReadInput() {
	// Setup fixture
	argsFixture := []string{"prog.name", "-i", "does.not.exist.txt"}

	// Exercise SUT
	actualCode := suite.sut.Run(argsFixture, nil, os.Stdout)

	// Verify results
	suite.Equal(cli.CouldNotReadInputCode, actualCode)
}

func (suite *EngineImplTestSuite) TestRun_GivenCouldNotOpenOutput_ShouldReturnCouldNotWriteOutput() {
	// Setup fixture
	// -> We make use of the fact that you cannot
	//    create a file with a name longer than 255 chars.
	argsFixture := []string{"prog.name", "-o", strings.Repeat("x", 2000)}

	// Exercise SUT
	actualCode := suite.sut.Run(argsFixture, os.Stdin, nil)

	// Verify results
	suite.Equal(cli.CouldNotWriteOutputCode, actualCode)
}

func (suite *EngineImplTestSuite) TestRun_GivenInvalidArgs_ShouldReturnFlagParseError() {
	// Setup fixture
	argsFixture := []string{"prog.name", "-b"}

	// Exercise SUT
	actualCode := suite.sut.Run(argsFixture, os.Stdin, nil)

	// Verify results
	suite.Equal(cli.FlagParseErrorCode, actualCode)
}
