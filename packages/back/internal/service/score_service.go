package service

import (
	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
	"github.com/saku0512/GYOUJI_HP/backend/internal/repository"
)

// テスト用インターフェース
type ScoreServiceInterface interface {
	GetScores() ([]model.ScoreBreakdown, error)
}

type ScoreService struct{ Repo *repository.ScoreRepository }

func NewScoreService(r *repository.ScoreRepository) *ScoreService { return &ScoreService{Repo: r} }

func (s *ScoreService) GetScores() ([]model.ScoreBreakdown, error) {
	return s.Repo.GetScores()
}
