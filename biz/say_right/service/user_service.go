package service

import (
	"context"
	"time"

	"api/biz/say_right/dal/model"
	"api/biz/say_right/dal/query"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	FindOrCreateUser(ctx context.Context, email string) (*model.User, error)
	SendVerificationCode(ctx context.Context, email string) error
	VerifyCode(ctx context.Context, email, code string) (bool, error)
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

func (s *userService) FindOrCreateUser(ctx context.Context, email string) (*model.User, error) {
	user, err := s.GetUserByEmail(ctx, email)
	if err == nil {
		return user, nil
	}
	// If not found, create
	newUser := &model.User{
		Email:     email,
		EmailNorm: email, // Assuming caller normalized it or we do it here
	}
	if err := s.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

func (s *userService) SendVerificationCode(ctx context.Context, email string) error {
	// MOCK: Just return success
	return nil
}

func (s *userService) VerifyCode(ctx context.Context, email, code string) (bool, error) {
	// MOCK: Check if code is 123456
	if code == "123456" {
		return true, nil
	}
	return false, nil
}
