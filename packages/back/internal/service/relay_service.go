package service

import (
	"errors"
	"sort"

	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
	"github.com/saku0512/GYOUJI_HP/backend/internal/repository"
)

type RelayService interface {
	GetRelayRankings(block string) (map[int]int, error)
	RegisterRelayRankings(block string, rankings map[int]int) error
}

type RelayServiceImpl struct {
	repo *repository.RelayRepository
}

func NewRelayService(repo *repository.RelayRepository) RelayService {
	return &RelayServiceImpl{repo: repo}
}

// GetRelayRankings 指定されたブロックのリレー順位を取得
func (s *RelayServiceImpl) GetRelayRankings(block string) (map[int]int, error) {
	// ブロック名のバリデーション
	if block != "A" && block != "B" {
		return nil, errors.New("block must be 'A' or 'B'")
	}

	return s.repo.GetRelayRankings(block)
}

// RegisterRelayRankings 指定されたブロックのリレー順位を登録
func (s *RelayServiceImpl) RegisterRelayRankings(block string, rankings map[int]int) error {
	// ブロック名のバリデーション
	if block != "A" && block != "B" {
		return errors.New("block must be 'A' or 'B'")
	}

	// 順位データのバリデーション
	if err := s.validateRankings(rankings); err != nil {
		return err
	}

	// 順位を登録
	if err := s.repo.RegisterRelayRankings(block, rankings); err != nil {
		return err
	}

	// 両ブロックの結果が揃った場合、ボーナス得点を計算
	completed, err := s.repo.CheckBothBlocksCompleted()
	if err != nil {
		return err
	}

	if completed {
		if err := s.calculateAndUpdateBonusScores(); err != nil {
			return err
		}
	}

	return nil
}

// validateRankings 順位データのバリデーション
func (s *RelayServiceImpl) validateRankings(rankings map[int]int) error {
	if len(rankings) != 6 {
		return errors.New("rankings must contain exactly 6 entries")
	}

	// 順位の重複チェック
	usedRanks := make(map[int]bool)
	usedGrades := make(map[int]bool)

	for rank, grade := range rankings {
		// 順位の範囲チェック
		if rank < 1 || rank > 6 {
			return errors.New("rank must be between 1 and 6")
		}

		// 学年の範囲チェック
		if grade < 1 || grade > 6 {
			return errors.New("grade must be between 1 and 6")
		}

		// 順位の重複チェック
		if usedRanks[rank] {
			return errors.New("duplicate rank found")
		}
		usedRanks[rank] = true

		// 学年の重複チェック
		if usedGrades[grade] {
			return errors.New("duplicate grade found")
		}
		usedGrades[grade] = true
	}

	// 1-6の順位が全て含まれているかチェック
	for i := 1; i <= 6; i++ {
		if !usedRanks[i] {
			return errors.New("missing rank: " + string(rune(i+'0')))
		}
	}

	return nil
}

// calculateAndUpdateBonusScores ボーナス得点を計算・更新
func (s *RelayServiceImpl) calculateAndUpdateBonusScores() error {
	// 全学年の得点を取得
	gradeScores, err := s.repo.GetGradeScores()
	if err != nil {
		return err
	}

	// 合計得点でソート
	sort.Slice(gradeScores, func(i, j int) bool {
		return gradeScores[i].TotalScore > gradeScores[j].TotalScore
	})

	// 最終順位とボーナス得点を計算
	s.calculateFinalRanksAndBonus(gradeScores)

	// ボーナス得点を更新
	return s.repo.UpdateRelayBonusScores(gradeScores)
}

// calculateFinalRanksAndBonus 最終順位とボーナス得点を計算
func (s *RelayServiceImpl) calculateFinalRanksAndBonus(gradeScores []model.GradeScore) {
	currentRank := 1
	
	for i := 0; i < len(gradeScores); i++ {
		// 同点の場合は同じ順位
		if i > 0 && gradeScores[i].TotalScore != gradeScores[i-1].TotalScore {
			currentRank = i + 1
		}
		
		gradeScores[i].FinalRank = currentRank
		
		// ボーナス得点を設定
		if bonus, exists := model.RelayBonusConfig[currentRank]; exists {
			gradeScores[i].BonusScore = bonus
		} else {
			gradeScores[i].BonusScore = 0
		}
	}
}