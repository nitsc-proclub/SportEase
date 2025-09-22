package service

import (
	"errors"

	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
	"github.com/saku0512/GYOUJI_HP/backend/internal/repository"
)

type AttendanceService interface {
	GetAttendanceScores() ([]model.AttendanceScore, error)
	UpdateAttendanceScore(classID int, score int) error
	BatchUpdateAttendanceScores(scores []model.AttendanceUpdateItem) error
}

type AttendanceServiceImpl struct {
	repo *repository.AttendanceRepository
}

func NewAttendanceService(repo *repository.AttendanceRepository) AttendanceService {
	return &AttendanceServiceImpl{repo: repo}
}

// GetAttendanceScores 全クラスの出席点を取得
func (s *AttendanceServiceImpl) GetAttendanceScores() ([]model.AttendanceScore, error) {
	return s.repo.GetAttendanceScores()
}

// UpdateAttendanceScore 出席点を更新
func (s *AttendanceServiceImpl) UpdateAttendanceScore(classID int, score int) error {
	// 入力値の検証
	if score < 0 || score > 10 {
		return errors.New("score must be between 0 and 10")
	}

	// 有効なクラスIDかチェック
	validClassIDs := []int{11, 12, 13, 21, 22, 23, 31, 32, 33, 41, 42, 43, 51, 52, 53, 6}
	isValid := false
	for _, id := range validClassIDs {
		if id == classID {
			isValid = true
			break
		}
	}
	if !isValid {
		return errors.New("invalid class_id")
	}

	return s.repo.UpdateAttendanceScore(classID, score)
}

// BatchUpdateAttendanceScores 出席点を一括更新
func (s *AttendanceServiceImpl) BatchUpdateAttendanceScores(scores []model.AttendanceUpdateItem) error {
	// 入力値の検証
	validClassIDs := map[int]bool{
		11: true, 12: true, 13: true,
		21: true, 22: true, 23: true,
		31: true, 32: true, 33: true,
		41: true, 42: true, 43: true,
		51: true, 52: true, 53: true,
		6: true,
	}

	for _, item := range scores {
		if item.Score < 0 || item.Score > 10 {
			return errors.New("all scores must be between 0 and 10")
		}
		if !validClassIDs[item.ClassID] {
			return errors.New("invalid class_id found")
		}
	}

	return s.repo.BatchUpdateAttendanceScores(scores)
}