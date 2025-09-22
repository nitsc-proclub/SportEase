package service

import (
	"github.com/saku0512/GYOUJI_HP/backend/internal/repository"
)

// テスト用インターフェース
type TournamentServiceInterface interface {
	GetTournamentsBySport(sport string) (interface{}, error)
}

type TournamentService struct {
	Repo *repository.TournamentRepository
}

func NewTournamentService(r *repository.TournamentRepository) *TournamentService {
	return &TournamentService{Repo: r}
}

func (s *TournamentService) GetTournamentsBySport(sport string) (interface{}, error) {
	return s.Repo.GetTournamentsBySport(sport)
}
