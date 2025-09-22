package repository

import (
	"database/sql"

	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
)

// UserRepository はユーザーデータへのアクセスに関するインターフェースです。
type UserRepository interface {
	FindUserByUsername(username string) (*model.User, error)
}

// userRepository は UserRepository の実装です。
type userRepository struct {
	db *sql.DB
}

// NewUserRepository は userRepository の新しいインスタンスを生成します。
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// FindUserByUsername はユーザー名でユーザーを検索します。
func (r *userRepository) FindUserByUsername(username string) (*model.User, error) {
	query := "SELECT id, username, hashed_password, role, assigned_sport FROM users WHERE username = ?"
	row := r.db.QueryRow(query, username)

	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.HashedPassword, &user.Role, &user.AssignedSport)
	if err != nil {
		// ユーザーが見つからない場合もエラーを返す
		return nil, err
	}

	return &user, nil
}
