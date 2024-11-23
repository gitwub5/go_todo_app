package service

import (
	"context"
	"fmt"

	"github.com/gitwub5/go_todo_app/store"
)

type Login struct {
	DB             store.Queryer
	Repo           UserGetter
	TokenGenerator TokenGenerator
}

func (l *Login) Login(ctx context.Context, name, pw string) (string, error) {
	// 사용자 정보 조회
	u, err := l.Repo.GetUser(ctx, l.DB, name)
	if err != nil {
		return "", fmt.Errorf("failed to list: %w", err)
	}
	// 비밀번호 비교
	if err := u.ComparePassword(pw); err != nil {
		return "", fmt.Errorf("wrong password: %w", err)
	}
	// JWT 생성
	jwt, err := l.TokenGenerator.GenerateToken(ctx, *u)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return string(jwt), nil
}
