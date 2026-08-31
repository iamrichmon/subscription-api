package service

import (
	"errors"
	"net/mail"
	"strconv"
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

var ErrEmailTaken = errors.New("Email already registered.")

func hashPlainPass(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func normalizeEmail(email string) (string, error) {

	addr, err := mail.ParseAddress(strings.TrimSpace(email))

	if err != nil {
		return "", err
	}

	return strings.ToLower(addr.Address), nil
}

func (s *UserService) Register(name, email, password string) (*model.User, error) {
	name = normalizeName(name)

	// name validation
	if name == "" || password == "" {
		return nil, errors.New("Name and password are required.")
	}

	// email validation

	addr, err := normalizeEmail(email)

	if err != nil {
		return nil, err
	}

	email = addr

	existing, err := s.repo.FindByEmail(email)
	if err == nil && existing != nil {
		return nil, ErrEmailTaken
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// password hashed

	hashed, err := hashPlainPass(password)

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

	addr, err := normalizeEmail(email)

	if err != nil {
		return "", err
	}

	email = addr

	existing, err := s.repo.FindByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err // DB/infra err, e.g. connection error
	}

	if existing == nil {
		return "", errors.New("invalid credentials") // no such user
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	//var jwtSecret = []byte("secret")

	return strconv.FormatUint(uint64(existing.ID), 10), nil

}
