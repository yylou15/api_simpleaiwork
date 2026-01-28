package service

import (
	"context"
	"time"

	"api/internal/dal/model"
	"api/internal/dal/query"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type userService struct {
	q *query.Query
}

func NewUserService() UserService {
	return &userService{
		q: query.Q,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) error {
	user.EmailVerifiedAt = time.Now().UTC()
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	return s.q.User.WithContext(ctx).Create(user)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.q.User.WithContext(ctx).Where(s.q.User.EmailNorm.Eq(email)).First()
}
