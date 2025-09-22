package service

import (
	"fmt"
	"strconv"

	"github.com/saku0512/GYOUJI_HP/backend/internal/repository"
)

// SettingService は設定に関するビジネスロジックのインターフェースです。
type SettingService interface {
	GetVisibility() (bool, error)
	UpdateVisibility(value bool) error
	GetWeather() (string, error)
	UpdateWeather(value string) error
}

// settingService は SettingService の実装です。
type settingService struct {
	repo repository.SettingRepository
}

// NewSettingService は settingService の新しいインスタンスを生成します。
func NewSettingService(repo repository.SettingRepository) SettingService {
	return &settingService{repo: repo}
}

// GetVisibility はスコア表示設定を取得します。
func (s *settingService) GetVisibility() (bool, error) {
	value, err := s.repo.GetSettingVisibility("showTotalScores")
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(value)
}

// UpdateVisibility はスコア表示設定を更新します。
func (s *settingService) UpdateVisibility(value bool) error {
	return s.repo.UpdateSettingVisibility("showTotalScores", strconv.FormatBool(value))
}

// GetWeather は天候設定を取得します。
func (s *settingService) GetWeather() (string, error) {
	return s.repo.GetSettingWeather("tableTennisWeather")
}

// UpdateWeather は天候設定を更新します。
func (s *settingService) UpdateWeather(value string) error {
	// 簡単なバリデーション
	if value != "sunny" && value != "rainy" {
		return fmt.Errorf("invalid weather value: %s", value)
	}
	return s.repo.UpdateSettingWeather("tableTennisWeather", value)
}
