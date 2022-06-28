package usecase

import "github.com/liampulles/span-digital-ranking-cli/pkg/league"

// Service provides usecases of the system, i.e. the real application logic.
type Service interface {
	CalculateRankings(gameResults []league.GameResult) []league.Ranking
}

type ServiceImpl struct{}

var _ Service = &ServiceImpl{}

func NewServiceImpl() *ServiceImpl {
	return &ServiceImpl{}
}

func (si *ServiceImpl) CalculateRankings(gameResults []league.GameResult) []league.Ranking {
	// Delegate to league package.
	return league.CalculateRankings(gameResults)
}
