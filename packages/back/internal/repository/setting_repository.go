package repository

import (
	"database/sql"
	"log"
)

// SettingRepository は設定データへのアクセスに関するインターフェースです。
type SettingRepository interface {
	GetSettingVisibility(key string) (string, error)
	UpdateSettingVisibility(key string, value string) error
	GetSettingWeather(key string) (string, error)
	UpdateSettingWeather(key string, value string) error
}

// settingRepository は SettingRepository の実装です。
type settingRepository struct {
	db *sql.DB
}

// NewSettingRepository は settingRepository の新しいインスタンスを生成します。
func NewSettingRepository(db *sql.DB) SettingRepository {
	return &settingRepository{db: db}
}

// GetSettingVisibility はキーで設定値を取得します。
func (r *settingRepository) GetSettingVisibility(key string) (string, error) {
	query := "SELECT setting_value FROM scoretable_settings WHERE setting_key = ?"
	row := r.db.QueryRow(query, key)

	var value string
	err := row.Scan(&value)
	if err != nil {
		return "", err
	}

	return value, nil
}

// UpdateSettingVisibility はキーで設定値を更新します。
func (r *settingRepository) UpdateSettingVisibility(key string, value string) error {
	query := "INSERT INTO scoretable_settings (setting_key, setting_value) VALUES (?, ?) ON DUPLICATE KEY UPDATE setting_value = ?"
	_, err := r.db.Exec(query, key, value, value)
	return err
}

func (r *settingRepository) GetSettingWeather(key string) (string, error) {
	query := "SELECT setting_value FROM weather_settings WHERE setting_key = ?"
	row := r.db.QueryRow(query, key)

	var value string
	err := row.Scan(&value)
	if err != nil {
		return "", err
	}

	log.Printf("[DEBUG] GetSettingWeather: key=%s, value=%s", key, value)

	return value, nil
}

func (r *settingRepository) UpdateSettingWeather(key string, value string) error {
	query := "INSERT INTO weather_settings (setting_key, setting_value) VALUES (?, ?) ON DUPLICATE KEY UPDATE setting_value = ?"
	_, err := r.db.Exec(query, key, value, value)
	return err
}