package league

import (
	"math"
	"sort"
)

// --- CalculateRankings related ---

type GameResult struct {
	TeamA  string
	ScoreA int
	TeamB  string
	ScoreB int
}

type Ranking struct {
	Rank   uint
	Team   string
	Points int
}

// Determine the ultimate ranking of all the teams in a league given game results.
func CalculateRankings(gameResults []GameResult) []Ranking {
	teamsToPoints := assignPointsForLeague(gameResults)
	return rankTeamPoints(teamsToPoints)
}

func assignPointsForLeague(gameResults []GameResult) map[string]int {
	teamsToPoints := make(map[string]int)
	for _, gameResult := range gameResults {
		pointsA, pointsB := AssignPoints(gameResult.ScoreA, gameResult.ScoreB)
		teamsToPoints[gameResult.TeamA] += pointsA
		teamsToPoints[gameResult.TeamB] += pointsB
	}
	return teamsToPoints
}

func rankTeamPoints(teamsToPoints map[string]int) []Ranking {
	// Just convert to a list
	rankings := make([]Ranking, len(teamsToPoints))
	count := -1
	for team, points := range teamsToPoints {
		count++
		rankings[count] = Ranking{Team: team, Points: points}
	}

	// Sort by points descending, then team name ascending
	sort.Slice(rankings, func(i int, j int) bool {
		a, b := rankings[i], rankings[j]
		if a.Points == b.Points {
			return a.Team < b.Team
		}
		return a.Points > b.Points
	})

	// Assign rank
	currRank := uint(0)
	currPoints := math.MaxInt
	for i, ranking := range rankings {
		if ranking.Points < currPoints {
			currRank = uint(i + 1)
			currPoints = ranking.Points
		}
		rankings[i].Rank = currRank
	}

	return rankings
}

// --- AssignPoints related ---

const (
	DrawPoints = 1
	WinPoints  = 3
	LosePoints = 0
)

// Assign points to A and B given their relative scores.
func AssignPoints(scoreA int, scoreB int) (pointsA int, pointsB int) {
	if scoreA == scoreB {
		return DrawPoints, DrawPoints
	}
	if scoreA > scoreB {
		return WinPoints, LosePoints
	}
	return LosePoints, WinPoints
}
