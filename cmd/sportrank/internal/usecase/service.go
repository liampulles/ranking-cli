package usecase

import "github.com/liampulles/span-digital-ranking-cli/pkg/league"

// Service provides usecases of the system, i.e. the real application logic.
type Service interface {
	CalculateRankings(gameResults []league.GameResult) []league.Ranking
}
