package league

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
	return nil
}

// --- Score related ---

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
