package services

import (
	"context"
	"errors"
	"strings"
	"user-srv/domain"
	"user-srv/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
