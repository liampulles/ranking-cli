package adapter

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/liampulles/ranking-cli/cmd/sportrank/internal/usecase"
	"github.com/liampulles/ranking-cli/pkg/league"
)

// Defined errors
var (
	ErrMalformedRow = errors.New("input row is malformed, it should be of the form <TeamA> <ScoreA>, <TeamB> <ScoreB>")
)

// RowIOGateway facilitates access to usecases of the system via "row"
// input and output.
type RowIOGateway interface {
	// Each input row should be of the form (ignoring quotes):
	// "<TeamA> <ScoreA>, <TeamB> <ScoreB>"
	// The resulting output rankings will be of the form (ignoring quotes):
	// "<Rank>. <Team>, <Points> <pt/pts>"
	CalculateRankings(rows []string) ([]string, error)
}

type RowIOGatewayImpl struct {
	usecaseSvc usecase.Service
}

var _ RowIOGateway = &RowIOGatewayImpl{}

func NewRowIOGatewayImpl(usecaseSvc usecase.Service) *RowIOGatewayImpl {
	return &RowIOGatewayImpl{
		usecaseSvc: usecaseSvc,
	}
}

func (riogi *RowIOGatewayImpl) CalculateRankings(rows []string) ([]string, error) {
	gameResults, err := riogi.convertInput(rows)
	if err != nil {
		return nil, err
	}

	rankings := riogi.usecaseSvc.CalculateRankings(gameResults)

	return riogi.convertOutput(rankings), nil
}

func (riogi *RowIOGatewayImpl) convertInput(rows []string) ([]league.GameResult, error) {
	gameResults := make([]league.GameResult, len(rows))
	for i, row := range rows {
		gameResult, err := riogi.convertInputRow(row)
		if err != nil {
			return nil, fmt.Errorf("could not convert row %d of input: %w", i, err)
		}

		gameResults[i] = gameResult
	}
	return gameResults, nil
}

const (
	rowSplitStr  = ","
	sideSplitStr = " "
)

func (riogi *RowIOGatewayImpl) convertInputRow(row string) (league.GameResult, error) {
	if strings.TrimSpace(row) == "" {
		return league.GameResult{}, fmt.Errorf("empty string: %w", ErrMalformedRow)
	}

	// Split into two sides, then parse each side.
	sides := strings.Split(row, rowSplitStr)
	if len(sides) != 2 {
		return league.GameResult{}, fmt.Errorf("expected 2 sections after splitting by comma but got %d: %w",
			len(sides), ErrMalformedRow)
	}

	teamA, scoreA, err := riogi.convertInputRowSide(sides[0])
	if err != nil {
		return league.GameResult{}, fmt.Errorf("first side: %w", err)
	}

	teamB, scoreB, err := riogi.convertInputRowSide(sides[1])
	if err != nil {
		return league.GameResult{}, fmt.Errorf("second side: %w", err)
	}

	return league.GameResult{
		TeamA:  teamA,
		ScoreA: scoreA,
		TeamB:  teamB,
		ScoreB: scoreB,
	}, nil
}

func (riogi *RowIOGatewayImpl) convertInputRowSide(side string) (string, int, error) {
	cleaned := strings.TrimSpace(side)

	// Since the name of the team may contain the split string, we only want to split on the LAST occurence.
	lastSpaceIdx := strings.LastIndex(cleaned, sideSplitStr)
	if lastSpaceIdx < 0 {
		return "", 0, fmt.Errorf("expected a space separating team and score but found none: %w", ErrMalformedRow)
	}

	team := strings.TrimSpace(cleaned[:lastSpaceIdx])
	scoreStr := strings.TrimSpace(cleaned[lastSpaceIdx+1:])

	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		return "", 0, fmt.Errorf("score is not an integer [%s]: %w", scoreStr, ErrMalformedRow)
	}

	return team, score, nil
}

func (riogi *RowIOGatewayImpl) convertOutput(rankings []league.Ranking) []string {
	rows := make([]string, len(rankings))
	for i, ranking := range rankings {
		rows[i] = riogi.convertOutputRanking(ranking)
	}
	return rows
}

func (riogi *RowIOGatewayImpl) convertOutputRanking(ranking league.Ranking) string {
	pointSuffix := riogi.determinePointSuffix(ranking.Points)
	return fmt.Sprintf("%d. %s, %d %s",
		ranking.Rank, ranking.Team, ranking.Points, pointSuffix)
}

const (
	pluralFormPointSuffix   = "pts"
	singularFormPointSuffix = "pt"
)

func (riogi *RowIOGatewayImpl) determinePointSuffix(points int) string {
	if points == 0 || points < -1 || points > 1 {
		return pluralFormPointSuffix
	}
	return singularFormPointSuffix
}
