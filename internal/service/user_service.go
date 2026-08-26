package service

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/iamrichmon/subscription-api/internal/model"
	"github.com/iamrichmon/subscription-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func normalizeName(name string) string {
	name = strings.TrimSpace(name)
	return strings.Join(strings.Fields(name), " ")
}

func (s *UserService) Register(name, email, password string) (*model.User, error) {
	name = normalizeName(name)

	// name validation
	if name == "" || password == "" {
		return nil, errors.New("Name and password are required.")
	}

	addr, err := mail.ParseAddress(strings.TrimSpace(email))

	if err != nil {
		return nil, err
	}

	// email validation

	email = strings.ToLower(addr.Address)

	existing, err := s.repo.FindByEmail(email)
	if err == nil && existing != nil {
		return nil, errors.New("Email already registered.")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// password hashed

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Plan:     model.FreeSubscription,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
