package services

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
	"user-srv/config"
	"user-srv/domain"
	"user-srv/repositories"

	"github.com/golang-jwt/jwt/v5"
)

type UserService interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error
	Login(ctx context.Context, email, password string) (string, error)
}

type userService struct {
	repo repositories.UserRepository
	cfg  *config.Config
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
		cfg:  config.LoadConfig(),
	}
}

func (s *userService) Create(ctx context.Context, user *domain.User) error {
	if strings.TrimSpace(user.Name) == "" {
		return errors.New("name cannot be empty")
	}
	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(user.Email, "@") {
		return errors.New("invalid email format")
	}
	if strings.TrimSpace(user.Password) == "" {
		return errors.New("password cannot be empty")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = hashedPassword

	return s.repo.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, id int) (*domain.User, error) {
	if id <= 0 {
		return nil, errors.New("id must be positive")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *userService) GetAll(ctx context.Context) ([]domain.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *userService) Update(ctx context.Context, user *domain.User) error {
	if user.ID <= 0 {
		return errors.New("id must be positive")
	}
	if strings.TrimSpace(user.Name) == "" {
		return errors.New("name cannot be empty")
	}
	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(user.Email, "@") {
		return errors.New("invalid email format")
	}
	if strings.TrimSpace(user.Password) == "" {
		return errors.New("password cannot be empty")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = hashedPassword

	return s.repo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("id must be positive")
	}
	return s.repo.Delete(ctx, id)
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	if strings.TrimSpace(email) == "" {
		return "", errors.New("email cannot be empty")
	}
	if strings.TrimSpace(password) == "" {
		return "", errors.New("password cannot be empty")
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
