package service

import (
	"errors"

	"strings"

	"github.com/iamrichmon/subscription-api/internal/auth"
	"github.com/iamrichmon/subscription-api/internal/model"
	"github.com/iamrichmon/subscription-api/internal/repository"
	"github.com/iamrichmon/subscription-api/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewUserService(repo *repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		repo:      repo,
		jwtSecret: jwtSecret}
}

func (s *UserService) Register(name, email, password string) (*model.User, error) {
	name = utils.NormalizeName(name)

	// name validation
	if name == "" || password == "" {
		return nil, errors.New("Name and password are required.")
	}

	// email validation

	addr, err := utils.NormalizeEmail(email)

	if err != nil {
		return nil, err
	}

	email = addr

	existing, err := s.repo.FindByEmail(email)
	if err == nil && existing != nil {
		return nil, utils.ErrEmailTaken
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// password hashed

	hashed, err := utils.HashPlainPass(password)

	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     name,
		Email:    email,
		Password: hashed,
		Plan:     model.FreeSubscription,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(email, password string) (string, error) {

	if strings.TrimSpace(password) == "" {
		return "", utils.ErrInvalidCredentials
	}

	addr, err := utils.NormalizeEmail(email)

	if err != nil {
		return "", err
	}

	email = addr

	existing, err := s.repo.FindByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err // DB/infra err, e.g. connection error
	}

	if existing == nil {
		return "", utils.ErrInvalidCredentials // no such user
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(password)); err != nil {
		return "", utils.ErrInvalidCredentials // password mismatch
	}

	token, err := auth.GenerateToken(existing.ID, existing.Email, s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
