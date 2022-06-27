package league

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

func CalculateRankings(gameResults []GameResult) []Ranking {
	return nil
}

func Score(scoreA int, scoreB int) (pointsA int, pointsB int) {
	return -1, -1
}
