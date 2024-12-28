package main

import (
	"game_rank/rank"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMain(m *testing.M) {
	gameInit()
}

func generateRankData() []*rank.RankSingleInfo {
	return []*rank.RankSingleInfo{
		{PlayerId: "player1", Score: 5000},
		{PlayerId: "player2", Score: 4800},
		{PlayerId: "player3", Score: 4600},
		{PlayerId: "player4", Score: 4500},
		{PlayerId: "player5", Score: 4300},
		{PlayerId: "player6", Score: 4100},
		{PlayerId: "player7", Score: 4000},
		{PlayerId: "player8", Score: 3800},
		{PlayerId: "player9", Score: 3600},
		{PlayerId: "player10", Score: 3400},
	}
}

func generatePlayerRankRangeData(playerId string) []*rank.RankSingleInfo {
	switch playerId {
	case "player1":
		return []*rank.RankSingleInfo{
			{PlayerId: "player1", Score: 5000},
			{PlayerId: "player2", Score: 4800},
			{PlayerId: "player3", Score: 4600},
			{PlayerId: "player4", Score: 4500},
			{PlayerId: "player5", Score: 4300},
		}
	case "player5":
		return []*rank.RankSingleInfo{
			{PlayerId: "player4", Score: 4500},
			{PlayerId: "player5", Score: 4300},
			{PlayerId: "player6", Score: 4100},
			{PlayerId: "player7", Score: 4000},
			{PlayerId: "player8", Score: 3800},
		}
	case "player10":
		return []*rank.RankSingleInfo{
			{PlayerId: "player8", Score: 3800},
			{PlayerId: "player9", Score: 3600},
			{PlayerId: "player10", Score: 3400},
		}
	default:
		return []*rank.RankSingleInfo{}
	}
}

func TestGetTopN(t *testing.T) {
	tests := []struct {
		name     string
		n        uint32
		expected []*rank.RankSingleInfo
		err      error
	}{
		{"Get top 5 players", 5, generateRankData()[:5], nil},
		{"Get top 10 players", 10, generateRankData(), nil},
		{"Invalid N value (zero)", 0, nil, errors.New("n is 0")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, topRanks := rank.GetTopN(tt.n)
			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expected, topRanks)
			}
		})
	}
}

// Test for fetching a specific player's rank and score
func TestGetPlayerRank(t *testing.T) {
	tests := []struct {
		name     string
		playerId string
		expected *rank.RankSingleInfo
		err      error
	}{
		{"Player1 exists", "player1", &rank.RankSingleInfo{PlayerId: "player1", Score: 5000}, nil},
		{"Player10 exists", "player10", &rank.RankSingleInfo{PlayerId: "player10", Score: 3400}, nil},
		{"Player does not exist", "player999", nil, errors.New("data is nil")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, rankInfo := rank.GetPlayerRank(tt.playerId)
			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expected, rankInfo)
			}
		})
	}
}

// Test for fetching player ranking range
func TestGetPlayerRankRange(t *testing.T) {
	tests := []struct {
		name     string
		playerId string
		expected []*rank.RankSingleInfo
	}{
		{"Player1 rank range", "player1", generatePlayerRankRangeData("player1")},
		{"Player5 rank range", "player5", generatePlayerRankRangeData("player5")},
		{"Player10 rank range", "player10", generatePlayerRankRangeData("player10")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, rankRange := rank.GetPlayerRankRange(tt.playerId, 5)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, rankRange)
		})
	}
}
