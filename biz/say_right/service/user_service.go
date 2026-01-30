package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"api/biz/say_right/dal/model"
	"api/biz/say_right/dal/query"
	"api/infra/mail"
	"api/infra/redis"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	FindOrCreateUser(ctx context.Context, email string) (*model.User, error)
	SendVerificationCode(ctx context.Context, email string) error
	VerifyCode(ctx context.Context, email, code string) (bool, error)
	UpgradeUserToPro(ctx context.Context, email string) error
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

const (
	CodeExpiration  = 5 * time.Minute
	KeyPrefixCode   = "biz:say_right:code:"   // code -> email (for uniqueness)
	KeyPrefixVerify = "biz:say_right:verify:" // email -> code (for verification)
	SendCooldown    = 60 * time.Second
	KeyPrefixSend   = "biz:say_right:send:" // email -> cooldown lock
)

type RateLimitError struct {
	RetryAfter int
}

func (e RateLimitError) Error() string {
	return fmt.Sprintf("Please wait %d seconds", e.RetryAfter)
}

func (s *userService) SendVerificationCode(ctx context.Context, email string) error {
	cooldownKey := KeyPrefixSend + email
	acquired, err := redis.Client.SetNX(ctx, cooldownKey, "1", SendCooldown).Result()
	if err != nil {
		return err
	}
	if !acquired {
		ttl, err := redis.Client.TTL(ctx, cooldownKey).Result()
		if err != nil {
			return err
		}
		retryAfter := int(ttl.Seconds())
		if retryAfter < 1 {
			retryAfter = int(SendCooldown.Seconds())
		}
		return RateLimitError{RetryAfter: retryAfter}
	}

	// 1. Clean up old code if exists
	verifyKey := KeyPrefixVerify + email
	oldCode, err := redis.Client.Get(ctx, verifyKey).Result()
	if err == nil && oldCode != "" {
		redis.Client.Del(ctx, KeyPrefixCode+oldCode)
	}

	// 2. Generate unique code
	var code string
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		c, err := generateCode()
		if err != nil {
			return err
		}

		// Check uniqueness
		exists, err := redis.Client.Exists(ctx, KeyPrefixCode+c).Result()
		if err != nil {
			return err
		}
		if exists == 0 {
			code = c
			break
		}
	}
	if code == "" {
		return errors.New("failed to generate unique code")
	}

	// 3. Save to Redis
	// Save code -> email (for uniqueness check)
	if err := redis.Client.Set(ctx, KeyPrefixCode+code, email, CodeExpiration).Err(); err != nil {
		return err
	}
	// Save email -> code (for verification)
	if err := redis.Client.Set(ctx, verifyKey, code, CodeExpiration).Err(); err != nil {
		return err
	}

	// 4. Send Email
	// Using empty subject to use default
	if err := mail.SendEmailCode(email, code, "Say Right Verify Code"); err != nil {
		return err
	}

	return nil
}

func (s *userService) VerifyCode(ctx context.Context, email, code string) (bool, error) {
	verifyKey := KeyPrefixVerify + email
	storedCode, err := redis.Client.Get(ctx, verifyKey).Result()
	if err != nil {
		// If key not found (expired or never sent)
		return false, nil
	}

	if storedCode == code {
		// Verification successful
		// Optional: Clean up after successful verification
		// redis.Client.Del(ctx, verifyKey)
		// redis.Client.Del(ctx, KeyPrefixCode+storedCode)
		return true, nil
	}

	return false, nil
}

func (s *userService) UpgradeUserToPro(ctx context.Context, email string) error {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	user.IsPro = true
	// Update specific columns using map since IsPro is not in generated query yet
	_, err = s.q.User.WithContext(ctx).Where(s.q.User.ID.Eq(user.ID)).Updates(map[string]interface{}{
		"is_pro": true,
	})
	return err
}

func generateCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
