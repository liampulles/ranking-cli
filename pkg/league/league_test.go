package league_test

import (
	"fmt"
	"testing"

	"github.com/liampulles/ranking-cli/pkg/league"
	"github.com/stretchr/testify/assert"
)

func TestCalculateRankings(t *testing.T) {
	// Setup fixture and expectations
	cases := []struct {
		gameResultsFixture []league.GameResult
		rankingsExpected   []league.Ranking
	}{
		// Trivial case
		{
			[]league.GameResult{},
			[]league.Ranking{},
		},

		// One game cases
		{
			[]league.GameResult{{"Albatros", 5, "Baboon", 2}},
			[]league.Ranking{
				{1, "Albatros", 3},
				{2, "Baboon", 0},
			},
		},
		{
			[]league.GameResult{{"Alphonse", 4, "Barry", 4}},
			[]league.Ranking{
				{1, "Alphonse", 1},
				{1, "Barry", 1},
			},
		},
		{
			[]league.GameResult{{"Barry", 4, "Alphonse", 4}},
			[]league.Ranking{
				{1, "Alphonse", 1},
				{1, "Barry", 1},
			},
		},

		// Mixed case
		{
			[]league.GameResult{
				{"Lions", 3, "Snakes", 3},
				{"Tarantulas", 1, "FC Awesome", 0},
				{"Lions", 1, "FC Awesome", 1},
				{"Tarantulas", 3, "Snakes", 1},
				{"Lions", 4, "Grouches", 0},
			},
			[]league.Ranking{
				{1, "Tarantulas", 6},
				{2, "Lions", 5},
				{3, "FC Awesome", 1},
				{3, "Snakes", 1},
				{5, "Grouches", 0},
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			// Exercise SUT
			rankingsActual := league.CalculateRankings(c.gameResultsFixture)

			// Verify results
			assert.Equal(t, c.rankingsExpected, rankingsActual)
		})
	}
}

func TestAssignPoints(t *testing.T) {
	// Setup fixture and expectations
	cases := []struct {
		scoreAFixture   int
		scoreBFixture   int
		pointsAExpected int
		pointsBExpected int
	}{
		// Draw cases
		{
			0, 0,
			1, 1,
		},
		{
			5, 5,
			1, 1,
		},

		// Win-lose cases
		{
			0, 1,
			0, 3,
		},
		{
			5, 2,
			3, 0,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			// Exercise SUT
			pointsAActual, pointsBActual := league.AssignPoints(c.scoreAFixture, c.scoreBFixture)

			// Verify results
			assert.Equal(t, c.pointsAExpected, pointsAActual)
			assert.Equal(t, c.pointsBExpected, pointsBActual)
		})
	}
}
