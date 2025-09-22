package service

import (
	"errors"
	"time"

	"github.com/saku0512/GYOUJI_HP/backend/internal/repository"

	"database/sql"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWTのClaims（ペイロード）の定義
type Claims struct {
	UserID        int64  `json:"user_id"`
	Username      string `json:"username"`
	Role          string `json:"role"`
	AssignedSport string `json:"assigned_sport"`
	jwt.RegisteredClaims
}

// UserService はユーザー関連のビジネスロジックのインターフェースです。
type UserService interface {
	Login(username, password string) (string, error)
}

type userService struct {
	userRepo  repository.UserRepository
	jwtSecret []byte
}

// NewUserService は userService の新しいインスタンスを生成します。
func NewUserService(repo repository.UserRepository, secret string) UserService {
	return &userService{
		userRepo:  repo,
		jwtSecret: []byte(secret),
	}
}

// Login はユーザーを認証し、JWTを返します。
func (s *userService) Login(username, password string) (string, error) {
	// 1. リポジトリを使ってユーザーを検索
	user, err := s.userRepo.FindUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			// ユーザーが存在しない
			return "", errors.New("invalid credentials")
		}
		// その他のデータベースエラー
		return "", err
	}

	// 2. パスワードを検証
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		// パスワードが一致しない
		return "", errors.New("invalid credentials")
	}

	// 3. JWTを生成
	expirationTime := time.Now().Add(24 * time.Hour) // トークンの有効期限を24時間に設定
	claims := &Claims{
		UserID:        user.ID,
		Username:      user.Username,
		Role:          user.Role,
		AssignedSport: user.AssignedSport,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
