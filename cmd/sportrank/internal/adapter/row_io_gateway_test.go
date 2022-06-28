package adapter_test

import (
	"fmt"
	"testing"

	"github.com/liampulles/ranking-cli/cmd/sportrank/internal/adapter"
	"github.com/liampulles/ranking-cli/cmd/sportrank/internal/usecase"
	"github.com/liampulles/ranking-cli/pkg/league"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RowIOGatewayImplTestSuite struct {
	suite.Suite
	mockUsecaseSvc *usecase.MockService
	sut            *adapter.RowIOGatewayImpl
}

func TestRowIOGatewayImplTestSuite(t *testing.T) {
	suite.Run(t, new(RowIOGatewayImplTestSuite))
}

func (suite *RowIOGatewayImplTestSuite) SetupTest() {
	suite.mockUsecaseSvc = usecase.NewMockService(suite.T())
	suite.sut = adapter.NewRowIOGatewayImpl(suite.mockUsecaseSvc)
}

func (suite *RowIOGatewayImplTestSuite) TestCalculateRankingsInput_InvalidCases() {
	// Setup fixture and expectations
	cases := []struct {
		fixture        []string
		expectedErrMsg string
	}{
		// Empty case
		{
			[]string{""},
			malformedRowErrMsg("could not convert row 0 of input: empty string"),
		},

		// "Side" split issues
		{
			[]string{"TeamA 1"},
			malformedRowErrMsg("could not convert row 0 of input: expected 2 sections after splitting by comma but got 1"),
		},
		{
			[]string{"TeamA 1, TeamB 2,"},
			malformedRowErrMsg("could not convert row 0 of input: expected 2 sections after splitting by comma but got 3"),
		},

		// "Side" issues
		{
			[]string{"TeamA, TeamB 1"},
			malformedRowErrMsg("could not convert row 0 of input: first side: expected a space separating team and score but found none"),
		},
		{
			[]string{"TeamA 1, TeamB seven"},
			malformedRowErrMsg("could not convert row 0 of input: second side: score is not an integer [seven]"),
		},

		// Error in one row of multiple
		{
			[]string{
				"TeamA 1, TeamB 2",
				"",
				"TeamA 3, TeamC 4",
			},
			malformedRowErrMsg("could not convert row 1 of input: empty string"),
		},
	}

	for i, c := range cases {
		suite.Run(fmt.Sprintf("Test case %d", i), func() {
			// Exercise SUT
			_, err := suite.sut.CalculateRankings(c.fixture)

			// Verify results
			suite.mockUsecaseSvc.AssertNotCalled(suite.T(), "CalculateRankings")
			suite.EqualError(err, c.expectedErrMsg)
		})
	}
}

func (suite *RowIOGatewayImplTestSuite) TestCalculateRankingsInput_ValidCases() {
	// Setup fixture
	cases := []struct {
		fixture            []string
		expectedConversion []league.GameResult
	}{
		// Trivial cases
		{nil, []league.GameResult{}},
		{[]string{}, []league.GameResult{}},

		// Mixed case - all rows should pass
		{
			[]string{
				"TeamA 1, TeamB 2",
				"John Lennon 7, Paul McCartney 2",
				" Dave Lister 1 , Arnold Rimmer 0 ",
			},
			[]league.GameResult{
				{"TeamA", 1, "TeamB", 2},
				{"John Lennon", 7, "Paul McCartney", 2},
				{"Dave Lister", 1, "Arnold Rimmer", 0},
			},
		},
	}

	for i, c := range cases {
		suite.Run(fmt.Sprintf("Test case %d", i), func() {
			// Setup mocks
			mockCall := suite.mockUsecaseSvc.Mock.
				On("CalculateRankings", c.expectedConversion).
				Return(nil)

			// Exercise SUT
			_, err := suite.sut.CalculateRankings(c.fixture)

			// Verify results
			suite.mockUsecaseSvc.AssertExpectations(suite.T())
			suite.NoError(err)

			// Cleanup
			mockCall.Unset()
		})
	}
}

func (suite *RowIOGatewayImplTestSuite) TestCalculateRankingsOutput() {
	// Setup fixture
	cases := []struct {
		mockOutput []league.Ranking
		expected   []string
	}{
		// Trivial cases
		{nil, []string{}},
		{[]league.Ranking{}, []string{}},

		// Mixed case - all rows should transform as expected
		{
			[]league.Ranking{
				{1, "John", 10},
				{1, "Paul", 10},
				{3, "George", 1},
				{4, "Ringo", 0},
				{5, "Peter", -1},
				{6, "Bob", -5},
			},
			[]string{
				"1. John, 10 pts",
				"1. Paul, 10 pts",
				"3. George, 1 pt",
				"4. Ringo, 0 pts",
				"5. Peter, -1 pt",
				"6. Bob, -5 pts",
			},
		},
	}

	for i, c := range cases {
		suite.Run(fmt.Sprintf("Test case %d", i), func() {
			// Setup mocks
			mockCall := suite.mockUsecaseSvc.Mock.
				On("CalculateRankings", mock.Anything).
				Return(c.mockOutput)

			// Exercise SUT
			actual, err := suite.sut.CalculateRankings(nil)

			// Verify results
			suite.NoError(err)
			suite.Equal(c.expected, actual)

			// Cleanup
			mockCall.Unset()
		})
	}
}

func malformedRowErrMsg(start string) string {
	return fmt.Sprintf("%s: %s", start, adapter.ErrMalformedRow.Error())
}
